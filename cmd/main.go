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
	"sync"
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
	errorBudget      int              // error budget. If exceeded, the sniper will stop and try to recover
)

const (
	DEFAULT_CHECK_INTERVAL = 300 // used if env CHECK_INTERVAL_SECONDS is not set or malformed
	MIN_CHECK_INTERVAL     = 60  // minimum check interval in seconds that makes sense

	ERROR_BUDGET_EXCEEDED_SLEEP = 10 * time.Second
	INITIAL_ERROR_BUDGET        = 5
	MAX_ERROR_BUDGET            = 10 // maximum error budget before being reset
	MIN_ERROR_BUDGET            = 1

	DEFAULT_NODE_DRAIN_TIMEOUT = 180              // used if env NODE_DRAIN_TIMEOUT_SECONDS is not set or malformed
	MIN_NODE_DRAIN_TIMEOUT     = 30               // minimum node drain timeout in seconds that makes sense
	NODE_DRAIN_SLEEP           = 10 * time.Second // sleep time after draining a node

	STATS_UPDATE_INTERVAL = 2 * time.Minute
)

func init() {

	healthy = true
	ready = true
	restoreErrorBudget()
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	var err error

	kubernetesClient, err = k8s.NewClient(nil)
	if !ok(err, logger, "failed to create Kubernetes client") {
		os.Exit(1)
	}

	googleClient, err = gcloud.NewClient(context.Background())
	if !ok(err, logger, "failed to create Google Cloud client") {
		os.Exit(2)
	}

	projectID, err = gcloud.GetProjectID()
	if !ok(err, logger, "failed to get project ID") {
		os.Exit(3)
	}

	allowedHours := os.Getenv("ALLOWED_HOURS")
	if allowedHours == "" {
		logger.Error("ALLOWED_HOURS environment variable is required")
		os.Exit(4)
	}
	allowedTimes, err = timing.ParseTimeSlots(strings.Split(allowedHours, ","))
	if !ok(err, logger, "failed to parse ALLOWED_HOURS") {
		os.Exit(5)
	}

	blockedHours := os.Getenv("BLOCKED_HOURS")
	if blockedHours != "" {
		blockedTimes, err = timing.ParseTimeSlots(strings.Split(blockedHours, ","))
		if !ok(err, logger, "failed to parse BLOCKED_HOURS") {
			os.Exit(6)
		}
	}

	checkIntervalStr := os.Getenv("CHECK_INTERVAL_SECONDS")
	if checkIntervalStr == "" {
		checkInterval = DEFAULT_CHECK_INTERVAL
	} else {
		checkInterval, err = strconv.Atoi(checkIntervalStr)
		if !ok(err, logger, "failed to parse CHECK_INTERVAL_SECONDS") {
			os.Exit(7)
		}
		if checkInterval <= MIN_CHECK_INTERVAL {
			checkInterval = DEFAULT_CHECK_INTERVAL
		}
	}

	nodeDrainTimeoutStr := os.Getenv("NODE_DRAIN_TIMEOUT_SECONDS")
	if nodeDrainTimeoutStr == "" {
		nodeDrainTimeout = DEFAULT_NODE_DRAIN_TIMEOUT
	} else {
		nodeDrainTimeout, err = strconv.Atoi(nodeDrainTimeoutStr)
		if !ok(err, logger, "failed to parse NODE_DRAIN_TIMEOUT_SECONDS") {
			os.Exit(8)
		}
	}
	if nodeDrainTimeout <= MIN_NODE_DRAIN_TIMEOUT {
		nodeDrainTimeout = DEFAULT_NODE_DRAIN_TIMEOUT
	}

	logger.Info("initialized", "project", projectID, "allowed", allowedTimes, "blocked", blockedTimes, "checkInterval", checkInterval, "nodeDrainTimeout", nodeDrainTimeout)
}

func main() {
	// Start web server for health checks and statistics
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
		err := http.ListenAndServe(":8080", nil)
		if !ok(err, logger, "failed to start HTTP server for health checks") {
			googleClient.Close()
			os.Exit(4)
		}
	}()

	// background goroutine for updating prometheus metrics
	go func() {
		for {
			stats.UpdateSnipedInLastHour()
			stats.UpdateSnipesExpectedInNextHour()
			time.Sleep(STATS_UPDATE_INTERVAL)
		}
	}()

	logger.Info("starting gke-preemptible-sniper")

	restoreErrorBudget()

	// main loop for checking nodes.
	for {
		checkErrorBudget(errorBudget)

		timeout := time.Duration(checkInterval) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		nodes, err := kubernetesClient.GetNodes(ctx)
		if !ok(err, logger, "failed to get nodes") {
			cancel()
			continue
		}
		logger.Info("retrieved nodes in the cluster", "amount", len(nodes))

		var wg sync.WaitGroup

		for _, node := range nodes {
			wg.Add(1)
			go func(node string) {
				nodeCtx, nodeCancel := context.WithTimeout(ctx, timeout)
				err := processNode(nodeCtx, node)

				if !ok(err, logger, "failed to process node", "node", node) {
					nodeCancel() // TODO check later if this is needed
					wg.Done()
					return
				}
				nodeCancel()
				wg.Done()
			}(node)
		}
		wg.Wait()
		cancel()
		increaseErrorBudget()

		logger.Info("sleeping", "seconds", checkInterval)
		time.Sleep(time.Duration(checkInterval) * time.Second)
	}
}

func processNode(ctx context.Context, node string) error {
	logger.Info("checking node", "node", node)
	hasAnnotation, err := kubernetesClient.HasNodeAnnotation(ctx, node, "gke-preemptible-sniper/timestamp")
	if !hasAnnotation {
		if !ok(err, logger, "failed to check sniper annotation", "error", err, "node", node) {
			return err
		}

		preemptibleAnnotation, err := kubernetesClient.HasNodeLabel(ctx, node, "cloud.google.com/gke-preemptible")
		if !ok(err, logger, "failed to check preemptible label", "error", err, "node", node) {
			return err
		}

		if !preemptibleAnnotation {
			logger.Info("skipping non-preemptible", "node", node)
			return nil
		}

		randTime, err := timing.CreateAllowedTime(allowedTimes, blockedTimes)
		if !ok(err, logger, "failed to create allowed time") {
			return err
		}

		logger.Info("adding annotation to node", "node", node, "timestamp", randTime.Format(time.RFC3339))
		err = kubernetesClient.SetNodeAnnotation(ctx, node, "gke-preemptible-sniper/timestamp", randTime.Format(time.RFC3339))
		if !ok(err, logger, "failed to add annotation", "error", err, "node", node) {
			return err
		}
	} else {
		logger.Info("node already has annotation", "node", node)
		// check if the node should be deleted
		timestamp, err := kubernetesClient.GetNodeAnnotation(ctx, node, "gke-preemptible-sniper/timestamp")
		if !ok(err, logger, "failed to get annotation", "error", err, "node", node) {
			return err
		}
		layout := time.RFC3339
		t, err := time.Parse(layout, timestamp)
		if !ok(err, logger, "failed to parse time", "error", err, "node", node, "timestamp", timestamp) {
			return err
		}
		if time.Now().After(t) {
			logger.Info("cordoning", "node", node)
			err = kubernetesClient.CordonNode(ctx, node)
			if !ok(err, logger, "failed to cordon node", "error", err, "node", node) {
				return err
			}
			drainCtx, drainCancel := context.WithTimeout(context.Background(), time.Duration(nodeDrainTimeout)*time.Second)
			logger.Info("draining", "node", node)
			err = kubernetesClient.DrainNode(drainCtx, node)
			if !ok(err, logger, "failed to drain node", "error", err, "node", node) {
				drainCancel()
				return err
			}
			drainCancel()
			time.Sleep(NODE_DRAIN_SLEEP)

			instance, err := kubernetesClient.GetNodeLabel(ctx, node, "kubernetes.io/hostname")
			if !ok(err, logger, "failed to get instance name", "error", err, "node", node) {
				return err
			}
			if instance == "" {
				logger.Error("instance name is empty", "node", node)
				return errors.New("instance name is empty")
			}

			zone, err := kubernetesClient.GetNodeZone(ctx, node)
			if !ok(err, logger, "failed to get zone", "error", err, "node", node) {
				return err
			}
			if zone == "" {
				logger.Error("zone is empty", "node", node)
				return errors.New("zone is empty")
			}

			logger.Info("deleting instance", "instance", instance, "zone", zone, "node", node)
			err = kubernetesClient.DeleteNode(ctx, node)
			if !ok(err, logger, "failed to delete node", "error", err, "node", node) {
				return err
			}

			err = googleClient.DeleteInstance(ctx, projectID, zone, instance)
			if !ok(err, logger, "failed to delete instance", "error", err, "instance", instance, "zone", zone, "project", projectID, "node", node) {
				return err
			}
			logger.Info("deleted instance", "instance", instance, "zone", zone, "project", projectID, "node", node)
			stats.AddSnipedNode(instance, time.Now())
		} else {
			duration := time.Until(t)
			logger.Info("node has time to live left", "node", node, "left", fmt.Sprintf("%vh%vm", int(duration.Hours()), int(duration.Minutes())%60))

			if duration < time.Hour {
				stats.AddExpectedSnipe(node, t)
			}
		}
	}
	return nil
}

func ok(err error, logger *slog.Logger, message string, loginfo ...any) bool {
	// add err to loginfo
	loginfo = append(loginfo, "error", err)
	if err != nil {
		logger.Error(message, loginfo...)
		decreaseErrorBudget()
		return false
	}
	return true
}

func checkErrorBudget(errorBudget int) {
	if errorBudget > MAX_ERROR_BUDGET {
		restoreErrorBudget()
	}
	if errorBudget <= MIN_ERROR_BUDGET {
		logger.Warn("error budget exceeded, trying to recover")
		healthy = false
		ready = false
		time.Sleep(ERROR_BUDGET_EXCEEDED_SLEEP)
		increaseErrorBudget()
	} else {
		healthy = true
		ready = true
	}
}

func decreaseErrorBudget() {
	errorBudget--
}

func increaseErrorBudget() {
	errorBudget++
}

func restoreErrorBudget() {
	errorBudget = INITIAL_ERROR_BUDGET
}
