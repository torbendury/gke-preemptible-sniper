package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

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

type Client struct {
	client kubernetes.Interface
}

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
	config.QPS = 100
	config.Burst = 500

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Client{client: clientset}, nil
}

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

func (c *Client) DrainNode(ctx context.Context, nodeName string) error {
	pods, err := c.client.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		if pod.Namespace == "kube-system" {
			continue // Skip system pods
		}

		err = c.evictPod(ctx, &pod)
		if err != nil {
			return err
		}
	}

	return nil
}

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
