package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/torbendury/gke-preemptible-sniper/gcloud"
	"github.com/torbendury/gke-preemptible-sniper/k8s"
	"github.com/torbendury/gke-preemptible-sniper/stats"
	"github.com/torbendury/gke-preemptible-sniper/timing"
)

var (
	allowedTimes     timing.TimeSlots // allowed times for node delete scheduling
	blockedTimes     timing.TimeSlots // blocked times for node delete scheduling
	checkInterval    int              // interval in seconds for checking nodes
	googleClient     *gcloud.Client   // Google Cloud client
	healthy          bool             // health status
	kubernetesClient *k8s.Client      // Kubernetes client
	logger           *slog.Logger     // logger
	nodeDrainTimeout int              // timeout in seconds for draining a node
	projectID        string           // Google Cloud project ID
	ready            bool             // readiness status

)

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

	projectID, err = gcloud.GetProjectID()
	if err != nil {
		logger.Error("failed to get project ID", "error", err)
		os.Exit(3)
	}

	allowedHours := os.Getenv("ALLOWED_HOURS")
	if allowedHours == "" {
		logger.Error("ALLOWED_HOURS environment variable is required")
		os.Exit(4)
	}
	allowedTimes, err = timing.ParseTimeSlots(strings.Split(allowedHours, ","))
	if err != nil {
		logger.Error("failed to parse ALLOWED_HOURS", "error", err)
		os.Exit(5)
	}

	blockedHours := os.Getenv("BLOCKED_HOURS")
	if blockedHours != "" {
		blockedTimes, err = timing.ParseTimeSlots(strings.Split(blockedHours, ","))
		if err != nil {
			logger.Error("failed to parse BLOCKED_HOURS", "error", err)
			os.Exit(6)
		}
	}

	checkIntervalStr := os.Getenv("CHECK_INTERVAL_SECONDS")
	if checkIntervalStr == "" {
		checkInterval = 1200
	} else {
		checkInterval, err = strconv.Atoi(checkIntervalStr)
		if err != nil {
			logger.Error("failed to parse CHECK_INTERVAL_SECONDS", "error", err)
			os.Exit(7)
		}
		if checkInterval == 0 {
			checkInterval = 1200
		}
	}

	nodeDrainTimeoutStr := os.Getenv("NODE_DRAIN_TIMEOUT_SECONDS")
	if nodeDrainTimeoutStr == "" {
		nodeDrainTimeout = 300
	} else {
		nodeDrainTimeout, err = strconv.Atoi(nodeDrainTimeoutStr)
		if err != nil {
			logger.Error("failed to parse NODE_DRAIN_TIMEOUT_SECONDS", "error", err)
			os.Exit(8)
		}
	}
	if nodeDrainTimeout <= 0 {
		nodeDrainTimeout = 300
	}

	healthy = true
	ready = true

	logger.Info("initialized", "project", projectID, "allowed", allowedTimes, "blocked", blockedTimes, "checkInterval", checkInterval, "nodeDrainTimeout", nodeDrainTimeout)
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

	http.Handle("/metrics", promhttp.HandlerFor(stats.Reg, promhttp.HandlerOpts{}))

	go func() {
		logger.Info("starting HTTP server for health checks")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			logger.Error("HTTP server error", "error", err)
			googleClient.Close()
			os.Exit(4)
		}
	}()

	go func() {
		for {
			logger.Info("updating sniped metrics")
			stats.UpdateSnipedInLastHour()
			stats.UpdateSnipesExpectedInNextHour()
			time.Sleep(2 * time.Minute)
		}
	}()

	logger.Info("starting gke-preemptible-sniper")

	errorBudget := 5

	// main loop
	for {

		checkErrorBudget(errorBudget)

		timeout := time.Duration(checkInterval) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		nodes, err := kubernetesClient.GetNodes(ctx)
		if err != nil {
			logger.Error("failed to get nodes", "error", err)
			cancel()
			errorBudget--
			continue
		}

		logger.Info("retrieved nodes in the cluster", "amount", len(nodes))
		for _, node := range nodes {
			err := processNode(ctx, node)
			if err != nil {
				logger.Error("failed to process node", "error", err, "node", node)
				errorBudget--
				continue
			}
		}
		cancel()
		errorBudget++

		logger.Info("sleeping", "seconds", checkInterval)
		time.Sleep(time.Duration(checkInterval) * time.Second)
	}
}

func checkErrorBudget(errorBudget int) {
	if errorBudget > 5 {
		errorBudget = 5
	}
	if errorBudget <= 0 {
		logger.Error("error budget exceeded, trying to recover")
		healthy = false
		ready = false
		time.Sleep(10 * time.Second)
		errorBudget++
	} else {
		healthy = true
		ready = true
	}
}

func processNode(ctx context.Context, node string) error {
	logger.Info("checking node", "node", node)
	hasAnnotation, err := kubernetesClient.HasNodeAnnotation(ctx, node, "gke-preemptible-sniper/timestamp")
	if !hasAnnotation {
		if err != nil {
			logger.Error("failed to check sniper annotation", "error", err, "node", node)
			return err
		}
		preemptibleAnnotation, err := kubernetesClient.HasNodeLabel(ctx, node, "cloud.google.com/gke-preemptible")
		if err != nil {
			logger.Error("failed to check preemptible label", "error", err, "node", node)
			return err
		}
		if !preemptibleAnnotation {
			logger.Info("skipping non-preemptible", "node", node)
			return nil
		}

		randTime, err := timing.CreateAllowedTime(allowedTimes, blockedTimes)
		if err != nil {
			logger.Error("failed to create allowed time", "error", err)
			return err
		}

		logger.Info("adding annotation to node", "node", node, "timestamp", randTime.Format(time.RFC3339))
		err = kubernetesClient.SetNodeAnnotation(ctx, node, "gke-preemptible-sniper/timestamp", randTime.Format(time.RFC3339))
		if err != nil {
			logger.Error("failed to add annotation", "error", err, "node", node)
			return err
		}
	} else {
		logger.Info("node already has annotation", "node", node)
		// check if the node should be deleted
		timestamp, err := kubernetesClient.GetNodeAnnotation(ctx, node, "gke-preemptible-sniper/timestamp")
		if err != nil {
			logger.Error("failed to check annotation", "error", err, "node", node)
			return err
		}
		layout := time.RFC3339
		t, err := time.Parse(layout, timestamp)
		if err != nil {
			logger.Error("failed to parse time", "error", err, "node", node, "timestamp", timestamp)
			return err
		}
		if time.Now().After(t) {
			logger.Info("cordoning", "node", node)
			err = kubernetesClient.CordonNode(ctx, node)
			if err != nil {
				logger.Error("failed to cordon node", "error", err, "node", node)
				return err
			}
			drainCtx, drainCancel := context.WithTimeout(context.Background(), time.Duration(nodeDrainTimeout)*time.Second)
			logger.Info("draining", "node", node)
			err = kubernetesClient.DrainNode(drainCtx, node)
			if err != nil {
				logger.Error("failed to drain node", "error", err, "node", node)
				drainCancel()
				return err
			}
			drainCancel()
			time.Sleep(10 * time.Second)

			instance, err := kubernetesClient.GetNodeLabel(ctx, node, "kubernetes.io/hostname")
			if err != nil {
				logger.Error("failed to get instance name", "error", err, "node", node)
			}
			if instance == "" {
				logger.Error("instance name is empty", "node", node)
				return errors.New("instance name is empty")
			}

			zone, err := kubernetesClient.GetNodeZone(ctx, node)
			if err != nil {
				logger.Error("failed to get zone", "error", err, "node", node)
			}
			if zone == "" {
				logger.Error("zone is empty", "node", node)
				return errors.New("zone is empty")
			}

			logger.Info("deleting instance", "instance", instance, "zone", zone, "node", node)
			err = googleClient.DeleteInstance(ctx, projectID, zone, instance)
			if err != nil {
				logger.Error("failed to delete instance", "error", err, "instance", instance, "zone", zone, "project", projectID, "node", node)
				return err
			}
			logger.Info("deleted instance", "instance", instance, "zone", zone, "project", projectID, "node", node)
			stats.AddSnipedNode(instance, time.Now())
		} else {
			duration := time.Until(t)
			logger.Info("node has time to live left", "node", node, "left", fmt.Sprintf("%vh%vm", int(duration.Hours()), int(duration.Minutes())%60))

			if duration < time.Hour {
				logger.Info("adding node to expected snipes metrics", "node", node, "timestamp", t)
				stats.AddExpectedSnipe(node, t)
			}
		}
	}
	return nil
}
