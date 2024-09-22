package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/torbendury/gke-preemptible-sniper/gcloud"
	"github.com/torbendury/gke-preemptible-sniper/k8s"
	"golang.org/x/exp/rand"
)

var kubernetesClient *k8s.Client
var googleClient *gcloud.Client
var projectID string

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	var err error
	kubernetesClient, err = k8s.NewClient(nil)
	if err != nil {
		logger.Error("failed to create Kubernetes client", "error", err)
		os.Exit(1)
	}

	googleClient, err = gcloud.NewClient(context.Background())
	if err != nil {
		logger.Error("failed to create Google Cloud client", "error", err)
		os.Exit(2)
	}

	// Get the project ID
	projectID, err = gcloud.GetProjectID()
	if err != nil {
		log.Fatalf("Failed to get project ID: %v", err)
	}
}

func getContextWithTimeout() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return ctx
}

func main() {

	logger.Info("starting gke-preemptible-sniper")

	// main loop
	for {
		logger.Debug("loop iteration for node check")
		ctx := getContextWithTimeout()

		nodes, err := kubernetesClient.GetNodes(ctx)
		if err != nil {
			logger.Error("failed to get nodes: %v", "error", err)
			logger.Info("retrying in 10 seconds")
			time.Sleep(10 * time.Second)
			continue
		}

		logger.Info("retrieved nodes in the cluster", "amount", len(nodes))
		for _, node := range nodes {
			logger.Info("checking node", "node", node)
			hasAnnotation, err := kubernetesClient.HasNodeAnnotation(ctx, node, "gke-preemptible-sniper/timestamp")
			if !hasAnnotation {
				if err != nil {
					logger.Error("failed to check annotation", "error", err, "node", node)
					continue
				}
				preemptibleAnnotation, err := kubernetesClient.GetNodeAnnotation(ctx, node, "cloud.google.com/gke-preemptible")
				if err != nil {
					logger.Error("failed to check preemptible annotation", "error", err, "node", node)
					continue
				}
				if preemptibleAnnotation != "true" {
					logger.Info("node is not preemptible", "node", node)
					continue
				}

				// calculate random time within 12h-18h in the future
				rand.Seed(uint64(time.Now().UnixNano()))
				randTime := time.Now().Add(time.Duration(rand.Intn(6)+12) * time.Hour)

				logger.Info("adding annotation to node", "node", node, "timestamp", randTime.Format(time.RFC3339))
				err = kubernetesClient.SetNodeAnnotation(ctx, node, "gke-preemptible-sniper/timestamp", randTime.Format(time.RFC3339))
				if err != nil {
					logger.Error("failed to add annotation", "error", err, "node", node)
					continue
				}
			} else {
				logger.Info("node already has annotation", "node", node)
				// check if the node should be deleted
				timestamp, err := kubernetesClient.GetNodeAnnotation(ctx, node, "gke-preemptible-sniper/timestamp")
				if err != nil {
					logger.Error("failed to check annotation", "error", err, "node", node)
					continue
				}
				layout := time.RFC3339
				t, err := time.Parse(layout, timestamp)
				if err != nil {
					logger.Error("failed to parse time", "error", err, "node", node, "timestamp", timestamp)
					continue
				}
				if time.Now().After(t) {
					logger.Info("node should be deleted", "node", node, "timestamp", t, "now", time.Now())
					// Cordon and Drain node, after that delete the GCP instance
					err = kubernetesClient.CordonNode(ctx, node)
					if err != nil {
						logger.Error("failed to cordon node", "error", err, "node", node)
						continue
					}
					err = kubernetesClient.DrainNode(ctx, node)
					if err != nil {
						logger.Error("failed to drain node", "error", err, "node", node)
						continue
					}
					// Get the instance name from the node
					instance, err := kubernetesClient.GetNodeLabel(ctx, node, "kubernetes.io/hostname")
					if err != nil {
						logger.Error("failed to get instance name", "error", err, "node", node)
					}

					// get the zone from the node
					zone, err := kubernetesClient.GetNodeZone(ctx, node)
					if err != nil {
						logger.Error("failed to get zone", "error", err, "node", node)
					}

					// Delete the instance
					err = googleClient.DeleteInstance(ctx, projectID, zone, instance)
					if err != nil {
						logger.Error("failed to delete instance", "error", err, "instance", instance, "zone", zone, "project", projectID, "node", node)
						continue
					}
				}
			}
		}
		time.Sleep(300 * time.Second)
	}
}
