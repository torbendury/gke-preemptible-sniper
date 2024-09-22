package k8s

import (
	"testing"

	"github.com/golang/mock/gomock"
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

	// Test with nil config to ensure fallback works
	client, err = NewClient(nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify that the client is not nil
	if client == nil {
		t.Fatalf("expected client, got nil")
	}
}
