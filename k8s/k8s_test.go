package k8s

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	testclient "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

// GetMockClient returns a fake k8s client for testing purposes
func GetMockClient() *Client {
	return &Client{
		client: testclient.NewSimpleClientset(),
	}
}

// MockKubernetesInterface is a mock of KubernetesInterface interface
type MockKubernetesInterface struct {
	kubernetes.Interface
}

func TestNewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock rest.Config
	mockConfig := &rest.Config{}

	// Inject the mock configuration into NewClient
	client, err := NewClient(mockConfig)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify that the client is not nil
	if client == nil {
		t.Fatalf("expected client, got nil")
	}

	// Test with nil config to ensure fallback also fails within CI
	client, err = NewClient(nil)
	if err == nil {
		t.Fatalf("expected error, got %v", err)
	}

	// Verify that the client is nil
	if client != nil {
		t.Fatalf("expected nil, got client")
	}
}

func TestGetNodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	client := GetMockClient()

	ctx := context.TODO()
	// Get the nodes
	nodes, err := client.GetNodes(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify that the nodes slice is empty
	if len(nodes) != 0 {
		t.Fatalf("expected no nodes, got %v", nodes)
	}
}

func TestCordonNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	client := GetMockClient()

	// Create mock node in client
	client.client.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node1",
		},
	}, metav1.CreateOptions{})

	// Cordon a node
	err := client.CordonNode(context.TODO(), "node1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDrainNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	client := GetMockClient()

	// Create mock node in client
	client.client.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node1",
		},
	}, metav1.CreateOptions{})

	// Drain a node
	err := client.DrainNode(context.TODO(), "node1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestEvictPod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	client := GetMockClient()

	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod1",
			Namespace: "default",
		},
	}

	// Create mock pod in client
	client.client.CoreV1().Pods("default").Create(context.TODO(), &pod, metav1.CreateOptions{})

	// Evict a pod
	err := client.evictPod(context.TODO(), &pod)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSetNodeAnnotation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	client := GetMockClient()

	// Create mock node in client
	client.client.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node1",
		},
	}, metav1.CreateOptions{})

	// Set node annotation
	err := client.SetNodeAnnotation(context.TODO(), "node1", "key", "value")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify that the node has the annotation
	node, err := client.client.CoreV1().Nodes().Get(context.TODO(), "node1", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	nodeValue := node.Annotations["key"]
	if nodeValue != "value" {
		t.Fatalf("expected value, got %v", nodeValue)
	}
}

func TestHasNodeAnnotation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	client := GetMockClient()

	// Create mock node in client
	client.client.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "node1",
			Annotations: map[string]string{"key": "value"},
		},
	}, metav1.CreateOptions{})

	// Check if node has annotation
	hasAnnotation, err := client.HasNodeAnnotation(context.TODO(), "node1", "key")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !hasAnnotation {
		t.Fatalf("expected annotation, got none")
	}
}

func TestGetNodeAnnotation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	client := GetMockClient()

	// Create mock node in client
	client.client.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "node1",
			Annotations: map[string]string{"key": "value"},
		},
	}, metav1.CreateOptions{})

	// Get node annotation
	value, err := client.GetNodeAnnotation(context.TODO(), "node1", "key")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if value != "value" {
		t.Fatalf("expected value, got %v", value)
	}
}

func TestGetNodeLabel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	client := GetMockClient()

	// Create mock node in client
	client.client.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node1",
			Labels: map[string]string{"key": "value"},
		},
	}, metav1.CreateOptions{})

	// Get node label
	value, err := client.GetNodeLabel(context.TODO(), "node1", "key")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if value != "value" {
		t.Fatalf("expected value, got %v", value)
	}
}

func TestGetNodeZone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	client := GetMockClient()

	// Create mock node in client
	client.client.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node1",
			Labels: map[string]string{"topology.kubernetes.io/zone": "zone1"},
		},
	}, metav1.CreateOptions{})

	client.client.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node2",
			Labels: map[string]string{"failure-domain.beta.kubernetes.io/zone": "zone1"},
		},
	}, metav1.CreateOptions{})

	// Get node zone
	zone, err := client.GetNodeZone(context.TODO(), "node1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if zone != "zone1" {
		t.Fatalf("expected zone1, got %v", zone)
	}

	zone, err = client.GetNodeZone(context.TODO(), "node2")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if zone != "zone1" {
		t.Fatalf("expected zone1, got %v", zone)
	}
}

func TestHasNodeLabel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	client := GetMockClient()

	// Create mock node in client
	client.client.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node1",
			Labels: map[string]string{"key": "value"},
		},
	}, metav1.CreateOptions{})

	// Check if node has label
	hasLabel, err := client.HasNodeLabel(context.TODO(), "node1", "key")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !hasLabel {
		t.Fatalf("expected label, got none")
	}
}
