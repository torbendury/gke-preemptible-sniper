package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/torbendury/gke-preemptible-sniper/gcloud"
	"github.com/torbendury/gke-preemptible-sniper/k8s"
	"golang.org/x/exp/rand"
)

func main() {
	client, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}
	googleClient, err := gcloud.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Google Cloud client: %v", err)
	}

	// Get the project ID
	projectID, err := gcloud.GetProjectID()
	if err != nil {
		log.Fatalf("Failed to get project ID: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// main loop
	for {

		nodes, err := client.GetNodes(ctx)
		if err != nil {
			log.Fatalf("Failed to get nodes: %v", err)
		}

		fmt.Printf("Nodes in the cluster: %d\n", len(nodes))
		for _, node := range nodes {
			fmt.Printf("Node: %v\n", node)
			hasAnnotation, err := client.HasNodeAnnotation(ctx, node, "gke-preemptible-sniper/timestamp")
			if !hasAnnotation {
				if err != nil {
					log.Fatalf("Failed to check annotation: %v", err)
				}

				// calculate random time within 12h-18h in the future
				rand.Seed(time.Now().UnixNano())
				randTime := time.Now().Add(time.Duration(rand.Intn(6)+12) * time.Hour)

				fmt.Printf("Adding annotation to node: %v\n", node)
				err = client.SetNodeAnnotation(ctx, node, "gke-preemptible-sniper/timestamp", randTime.Format(time.RFC3339))
				if err != nil {
					log.Fatalf("Failed to add annotation: %v", err)
				}
			} else {
				fmt.Printf("Node %v already has annotation\n", node)
				// check if the node should be deleted
				timestamp, err := client.GetNodeAnnotation(ctx, node, "gke-preemptible-sniper/timestamp")
				if err != nil {
					log.Fatalf("Failed to get annotation: %v", err)
				}
				layout := time.RFC3339
				t, err = time.Parse(layout, timestamp)
				if time.Now().After(t) {
					fmt.Printf("Node %v should be deleted\n", node)
					// Cordon and Drain node, after that delete the GCP instance
					err = client.CordonNode(ctx, node)
					if err != nil {
						log.Fatalf("Failed to cordon node: %v", err)
					}
					err = client.DrainNode(ctx, node)
					if err != nil {
						log.Fatalf("Failed to drain node: %v", err)
					}
					// Get the instance name from the node
					instance, err := client.GetNodeLabel(ctx, node, "kubernetes.io/hostname")
					if err != nil {
						log.Fatalf("Failed to get instance name: %v", err)
					}

					// get the zone from the node
					zone, err := client.GetNodeZone(ctx, node)
					if err != nil {
						log.Fatalf("Failed to get zone: %v", err)
					}

					// Delete the instance
					err = googleClient.DeleteInstance(ctx, projectID, zone, instance)
				}
			}
		}

		time.Sleep(300 * time.Second)
	}
}
