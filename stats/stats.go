// Package stats provides functions for calculating and storing internal Prometheus metrics about sniped GKE nodes.
package stats

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type SnipedNode struct {
	NodeName string
	Time     time.Time
}

// SnipedNodes is a list of sniped nodes.
type SnipedNodes []SnipedNode

var (
	SnipedInLastHour = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gke_preemptible_sniper_sniped_last_hour",
		Help: "Number of nodes sniped in the last hour",
	}, []string{"node", "time"})

	SnipesExpectedInNextHour = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gke_preemptible_sniper_snipes_expected_next_hour",
		Help: "Number of nodes expected to be sniped in the next hour",
	}, []string{"node", "time"})

	Reg = prometheus.NewRegistry()

	snipedInLastHour         SnipedNodes
	snipesExpectedInNextHour SnipedNodes
)

func init() {
	Reg.MustRegister(SnipedInLastHour, SnipesExpectedInNextHour)

	snipedInLastHour = make(SnipedNodes, 0)
}

// AddSnipedNode adds a sniped node to the list of sniped nodes.
func AddSnipedNode(nodeName string, time time.Time) {
	snipedInLastHour = append(snipedInLastHour, SnipedNode{NodeName: nodeName, Time: time})
}

// UpdateSnipedInLastHour updates the number of sniped nodes in the last hour. It removes nodes that are older than an hour.
func UpdateSnipedInLastHour() {
	var snipedNodes SnipedNodes
	for _, snipedNode := range snipedInLastHour {
		if snipedNode.Time.After(time.Now().Add(-time.Hour)) {
			snipedNodes = append(snipedNodes, snipedNode)
		}
	}

	snipedInLastHour = snipedNodes

	for _, snipedNode := range snipedInLastHour {
		SnipedInLastHour.WithLabelValues(snipedNode.NodeName, snipedNode.Time.Format(time.RFC3339)).Set(1)
	}
}

// AddExpectedSnipe adds an expected snipe to the list of expected snipes.
func AddExpectedSnipe(nodeName string, time time.Time) {
	snipesExpectedInNextHour = append(snipesExpectedInNextHour, SnipedNode{NodeName: nodeName, Time: time})
}

// UpdateSnipesExpectedInNextHour updates the number of expected snipes in the next hour. It removes nodes which timestamp has already passed or is further than an hour away.
func UpdateSnipesExpectedInNextHour() {
	var snipedNodes SnipedNodes
	for _, snipedNode := range snipesExpectedInNextHour {
		if snipedNode.Time.After(time.Now()) && snipedNode.Time.Before(time.Now().Add(time.Hour)) {
			snipedNodes = append(snipedNodes, snipedNode)
		}
	}

	snipesExpectedInNextHour = snipedNodes

	for _, snipedNode := range snipesExpectedInNextHour {
		SnipesExpectedInNextHour.WithLabelValues(snipedNode.NodeName, snipedNode.Time.Format(time.RFC3339)).Set(1)
	}
}
