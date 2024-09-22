package gcloud

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
)

func TestGetProjectID(t *testing.T) {
	expectedProjectID := "test-project-id"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/computeMetadata/v1/project/project-id" {
			t.Fatalf("unexpected URL path: %s", r.URL.Path)
		}
		if r.Header.Get("Metadata-Flavor") != "Google" {
			t.Fatalf("unexpected Metadata-Flavor header: %s", r.Header.Get("Metadata-Flavor"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedProjectID))
	}))
	defer server.Close()

	originalMetadataURL := metadataURL
	metadataURL = server.URL + "/computeMetadata/v1/project/project-id"
	defer func() { metadataURL = originalMetadataURL }()

	projectID, err := GetProjectID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if projectID != expectedProjectID {
		t.Fatalf("expected %s, got %s", expectedProjectID, projectID)
	}
}

func TestListInstances(t *testing.T) {
	mockClient := &mockInstancesClient{
		listFunc: func(ctx context.Context, req *computepb.ListInstancesRequest, opts ...gax.CallOption) *mockInstanceIterator {
			return &mockInstanceIterator{
				items: []*computepb.Instance{
					{Name: stringPtr("instance-1")},
					{Name: stringPtr("instance-2")},
				},
				nextFunc: func() error {
					return nil
				},
			}
		},
	}

	ctx := context.Background()
	req := &computepb.ListInstancesRequest{}
	it := mockClient.List(ctx, req)

	var instances []*computepb.Instance
	for {
		instance, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		instances = append(instances, instance)
	}

	if len(instances) != 2 {
		t.Fatalf("expected 2 instances, got %d", len(instances))
	}
}

// Helper function to create a pointer to a string
func stringPtr(s string) *string {
	return &s
}

type mockInstancesClient struct {
	listFunc func(ctx context.Context, req *computepb.ListInstancesRequest, opts ...gax.CallOption) *mockInstanceIterator
}

func (m *mockInstancesClient) List(ctx context.Context, req *computepb.ListInstancesRequest, opts ...gax.CallOption) *mockInstanceIterator {
	return m.listFunc(ctx, req, opts...)
}

type mockInstanceIterator struct {
	items    []*computepb.Instance
	pageInfo *iterator.PageInfo
	nextFunc func() error
}

func (it *mockInstanceIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

func (it *mockInstanceIterator) Next() (*computepb.Instance, error) {
	if len(it.items) == 0 {
		if err := it.nextFunc(); err != nil {
			return nil, err
		}
	}
	if len(it.items) == 0 {
		return nil, iterator.Done
	}
	var item *computepb.Instance
	item, it.items = it.items[0], it.items[1:]
	return item, nil
}
