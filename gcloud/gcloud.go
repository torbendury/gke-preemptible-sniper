package gcloud

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
)

type Client struct {
	client *compute.InstancesClient
}

var metadataURL = "http://metadata.google.internal/computeMetadata/v1/project/project-id"

// NewClient creates a new Google Cloud client using the provided context (workload identity) and returns a Client.
func NewClient(ctx context.Context) (*Client, error) {
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return &Client{client: client}, nil
}

// GetProjectID retrieves the project ID from the metadata server.
func GetProjectID() (string, error) {
	metadataURL := metadataURL
	req, err := http.NewRequest("GET", metadataURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Metadata-Flavor", "Google")

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get project ID, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) ListInstances(ctx context.Context, projectID, zone string) ([]*computepb.Instance, error) {
	req := &computepb.ListInstancesRequest{
		Project: projectID,
		Zone:    zone,
	}

	var instances []*computepb.Instance
	it := c.client.List(ctx, req)
	for {
		instance, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list instances: %v", err)
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (c *Client) DeleteInstance(ctx context.Context, projectID, zone, instanceName string) error {
	req := &computepb.DeleteInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	op, err := c.client.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete instance: %v", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for the delete operation: %v", err)
	}

	fmt.Printf("Instance %s deleted successfully\n", instanceName)
	return nil
}

// GetInstance retrieves a single instance in the specified project, zone, and instance name.
func (c *Client) GetInstance(ctx context.Context, projectID, zone, instanceName string) (*computepb.Instance, error) {
	req := &computepb.GetInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	instance, err := c.client.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %v", err)
	}

	return instance, nil
}
