package gcloud

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type Client struct {
	computeService *compute.Service
}

// NewClient creates a new Client with a Google Compute Engine service using default credentials.
func NewClient(ctx context.Context) (*Client, error) {
	// Use default credentials
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile("/var/run/secrets/google/key.json"))
	if err != nil {
		return nil, err
	}

	return &Client{computeService: computeService}, nil
}

// GetProjectID retrieves the project ID from the metadata server.
func GetProjectID() (string, error) {
	const metadataURL = "http://metadata.google.internal/computeMetadata/v1/project/project-id"
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

// ListInstances lists all instances in the specified project and zone.
func (c *Client) ListInstances(ctx context.Context, projectID, zone string) ([]*compute.Instance, error) {
	req := c.computeService.Instances.List(projectID, zone)
	var instances []*compute.Instance
	if err := req.Pages(ctx, func(page *compute.InstanceList) error {
		instances = append(instances, page.Items...)
		return nil
	}); err != nil {
		return nil, err
	}
	return instances, nil
}

// GetInstance retrieves a single instance in the specified project, zone, and instance name.
func (c *Client) GetInstance(ctx context.Context, projectID, zone, instanceName string) (*compute.Instance, error) {
	instance, err := c.computeService.Instances.Get(projectID, zone, instanceName).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// DeleteInstance deletes an instance in the specified project, zone, and instance name.
func (c *Client) DeleteInstance(ctx context.Context, projectID, zone, instanceName string) error {
	_, err := c.computeService.Instances.Delete(projectID, zone, instanceName).Context(ctx).Do()
	if err != nil {
		return err
	}
	return nil
}
