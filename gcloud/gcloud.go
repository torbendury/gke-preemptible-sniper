// Package gcloud provides a client for interacting with Google Cloud Platform.
// Is provides methods for listing, deleting, and getting instances, as well as authenticating with the Google Cloud API and retrieving the GCP project ID.
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

const (
	MAXIMUM_RETRIES = 3
)

// Client is a Google Cloud client.
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
	return &Client{client: client}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
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

	var lastErr error
	for i := 0; i < MAXIMUM_RETRIES; i++ {
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("failed to get project ID, status code: %d", resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		return string(body), nil
	}
	return "", lastErr
}

// ListInstances retrieves a list of instances in the specified project and zone.
func (c *Client) ListInstances(ctx context.Context, projectID, zone string) ([]*computepb.Instance, error) {
	req := &computepb.ListInstancesRequest{
		Project: projectID,
		Zone:    zone,
	}

	var lastErr error
	var thisErr error

	for i := 0; i < MAXIMUM_RETRIES; i++ {
		var instances []*computepb.Instance
		it := c.client.List(ctx, req)
		for {
			lastErr = thisErr
			thisErr = nil
			instance, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				thisErr = fmt.Errorf("failed to list instances: %v", err)
				break
			}
			instances = append(instances, instance)
		}
		if thisErr == nil {
			continue
		}
		return instances, nil
	}
	return nil, lastErr
}

// DeleteInstance deletes an instance in the specified project, zone, and instance name.
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

	var lastErr error
	for i := 0; i < MAXIMUM_RETRIES; i++ {

		instance, err := c.client.Get(ctx, req)
		if err != nil {
			lastErr = fmt.Errorf("failed to get instance: %v", err)
			continue
		}
		return instance, nil
	}
	return nil, lastErr
}
