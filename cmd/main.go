package main

import (
	"context"
	"log/slog"
	"net/http"
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

var healthy bool
var ready bool

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

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
		logger.Error("failed to get project ID", "error", err)
		os.Exit(3)
	}

	healthy = true
	ready = true
}

func getContextWithTimeout() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx, cancel
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if healthy {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("not ok"))
		}
	})

	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if ready {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ready"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("not ready"))
		}
	})

	go func() {
		logger.Info("starting HTTP server for health checks")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			logger.Error("HTTP server error", "error", err)
			os.Exit(4)
		}
	}()

	logger.Info("starting gke-preemptible-sniper")

	// main loop
	for {
		logger.Debug("loop iteration for node check")
		ctx, cancel := getContextWithTimeout()

		nodes, err := kubernetesClient.GetNodes(ctx)
		if err != nil {
			logger.Error("failed to get nodes", "error", err)
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
					logger.Error("failed to check sniper annotation", "error", err, "node", node)
					continue
				}
				preemptibleAnnotation, err := kubernetesClient.HasNodeLabel(ctx, node, "cloud.google.com/gke-preemptible")
				if err != nil {
					logger.Error("failed to check preemptible label", "error", err, "node", node)
					continue
				}
				if !preemptibleAnnotation {
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
		cancel()
		time.Sleep(300 * time.Second)
	}
}
