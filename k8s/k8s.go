// Package k8s provides a client for interacting with Kubernetes.
// It provides methods for listing, cordoning, draining, and deleting nodes, as well as setting and getting node annotations and labels.
package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Client is a Kubernetes client wrapper.
type Client struct {
	client kubernetes.Interface
}

type PodEvictionError struct {
	PodName      string
	PodNamespace string
	Err          error
}

func (r *PodEvictionError) Error() string {
	return fmt.Sprintf("pod %v, namespace %v, err %v", r.PodName, r.PodNamespace, r.Err)
}

const POD_EVICT_TIMEOUT_SECONDS = 30

// NewClient creates a new Kubernetes client using the provided rest.Config and returns a Client.
// If no config is provided, it will try to use in-cluster config, and if that fails, it will fallback to a kubeconfig file.
// At the moment, it does not apply client-side rate limiting.
func NewClient(config *rest.Config) (*Client, error) {
	var err error

	if config == nil {
		// Try to use in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			// Fallback to kubeconfig file
			var kubeconfig string
			if home := homedir.HomeDir(); home != "" {
				kubeconfig = filepath.Join(home, ".kube", "config")
			} else {
				kubeconfig = ""
			}

			config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				return nil, err
			}
		}
	}

	// Disable client request limiting
	config.QPS = -1
	config.Burst = -1

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Client{client: clientset}, nil
}

// GetNodes returns a list of node names in the Kubernetes cluster where the client points to.
func (c *Client) GetNodes(ctx context.Context) ([]string, error) {
	nodes, err := c.client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeNames []string
	for _, node := range nodes.Items {
		nodeNames = append(nodeNames, node.Name)
	}

	return nodeNames, nil
}

// CordonNode cordon the node with the provided name.
// Since no out of the box method is provided by the client-go library, we need to patch the node object to set the spec.unschedulable field to true.
func (c *Client) CordonNode(ctx context.Context, nodeName string) error {
	node, err := c.client.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Create a patch to update the node's spec.unschedulable field
	oldData, err := json.Marshal(node)
	if err != nil {
		return err
	}

	node.Spec.Unschedulable = true

	newData, err := json.Marshal(node)
	if err != nil {
		return err
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, node)
	if err != nil {
		return err
	}

	_, err = c.client.CoreV1().Nodes().Patch(ctx, nodeName, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
	return err
}

// DrainNode drains the node with the provided name.
// It evicts all the pods running on the node, except for the ones in the kube-system namespace.
// It uses the Eviction API to evict the pods.
func (c *Client) DrainNode(ctx context.Context, nodeName string) error {
	pods, err := c.client.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(pods.Items))

	for _, pod := range pods.Items {
		if pod.Namespace == "kube-system" {
			continue // Skip system pods
		}

		wg.Add(1)

		go func(pod v1.Pod) {

			defer wg.Done()

			err = c.evictPod(ctx, &pod)
			if err != nil {
				// try to recover by deleting the pod
				err = c.DeletePod(ctx, pod.Name, pod.Namespace)
				if err != nil {
					errChan <- err
				}
			}

			for range POD_EVICT_TIMEOUT_SECONDS {
				_, err = c.client.CoreV1().Pods(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
				if err != nil {
					break
				}
				<-time.After(1 * time.Second)
			}
			_, err = c.client.CoreV1().Pods(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
			if err == nil {
				errChan <- &PodEvictionError{
					PodName:      pod.Name,
					PodNamespace: pod.Namespace,
					Err:          fmt.Errorf("pod %s/%s still exists after eviction", pod.Namespace, pod.Name),
				}
			}

		}(pod)
	}

	wg.Wait()
	close(errChan)

	// Get all errors and compile them to one
	var errList []error
	for err := range errChan {
		errList = append(errList, err)
	}
	if len(errList) > 0 {
		return fmt.Errorf("failed to drain node %s: %v", nodeName, errList)
	}

	return nil
}

// evictPod evicts the provided pod.
func (c *Client) evictPod(ctx context.Context, pod *v1.Pod) error {
	eviction := &v1beta1.Eviction{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
		DeleteOptions: &metav1.DeleteOptions{},
	}

	return c.client.CoreV1().Pods(pod.Namespace).Evict(ctx, eviction)
}

// SetNodeAnnotation sets the provided key-value pair as an annotation on the node with the provided name.
// It uses the Patch API to update the node's annotations.
func (c *Client) SetNodeAnnotation(ctx context.Context, nodeName, key, value string) error {
	node, err := c.client.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Create a patch to update the node's annotations
	oldData, err := json.Marshal(node)
	if err != nil {
		return err
	}

	if node.Annotations == nil {
		node.Annotations = make(map[string]string)
	}
	node.Annotations[key] = value

	newData, err := json.Marshal(node)
	if err != nil {
		return err
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, node)
	if err != nil {
		return err
	}

	_, err = c.client.CoreV1().Nodes().Patch(ctx, nodeName, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
	return err
}

// HasNodeAnnotation checks if the node with the provided name has the annotation with the provided key.
func (c *Client) HasNodeAnnotation(ctx context.Context, nodeName, key string) (bool, error) {
	node, err := c.client.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	if node.Annotations == nil {
		return false, nil
	}

	_, exists := node.Annotations[key]
	return exists, nil
}

// GetNodeAnnotation returns the value of the annotation with the provided key on the node with the provided name.
func (c *Client) GetNodeAnnotation(ctx context.Context, nodeName, key string) (string, error) {
	node, err := c.client.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if node.Annotations == nil {
		return "", fmt.Errorf("no annotations found on node %s", nodeName)
	}

	value, exists := node.Annotations[key]
	if !exists {
		return "", fmt.Errorf("annotation %s not found on node %s", key, nodeName)
	}

	return value, nil
}

// GetNodeLabel returns the value of the label with the provided key on the node with the provided name.
func (c *Client) GetNodeLabel(ctx context.Context, nodeName, key string) (string, error) {
	node, err := c.client.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if node.Labels == nil {
		return "", fmt.Errorf("no labels found on node %s", nodeName)
	}

	value, exists := node.Labels[key]
	if !exists {
		return "", fmt.Errorf("label %s not found on node %s", key, nodeName)
	}

	return value, nil
}

// HasNodeLabel checks if the node with the provided name has the label with the provided key.
func (c *Client) HasNodeLabel(ctx context.Context, nodeName, key string) (bool, error) {
	node, err := c.client.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	if node.Labels == nil {
		return false, nil
	}

	_, exists := node.Labels[key]
	return exists, nil
}

// GetNodeZone returns the zone of the node with the provided name.
// It first tries to get the zone from the "topology.kubernetes.io/zone" label, and if that fails, it tries to get it from the "failure-domain.beta.kubernetes.io/zone" label.
func (c *Client) GetNodeZone(ctx context.Context, nodeName string) (string, error) {
	l, err := c.GetNodeLabel(ctx, nodeName, "topology.kubernetes.io/zone")
	if err != nil {
		l, err = c.GetNodeLabel(ctx, nodeName, "failure-domain.beta.kubernetes.io/zone")
		if err != nil {
			return "", fmt.Errorf("zone label not found on node %s", nodeName)
		}
	}
	return l, nil
}

// DeleteNode deletes the node with the provided name.
func (c *Client) DeleteNode(ctx context.Context, nodeName string) error {
	return c.client.CoreV1().Nodes().Delete(ctx, nodeName, metav1.DeleteOptions{})
}

// DeletePod deletes the pod with the provided name and namespace.
func (c *Client) DeletePod(ctx context.Context, podName, namespace string) error {
	return c.client.CoreV1().Pods(namespace).Delete(ctx, podName, metav1.DeleteOptions{})
}

// GetPod checks if a Pod exists in the given namespace and returns it.
func (c *Client) GetPod(ctx context.Context, podName, namespace string) (*v1.Pod, error) {
	pod, err := c.client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}
