# gke-preemptible-sniper

`gke-preemptible-sniper` is a small application written in Golang that is supposed to run inside a Google Kubernetes cluster and work around the [known limitation of preemptible VMs](https://cloud.google.com/compute/docs/instances/preemptible#limitations).

Its' purpose is to gracefully remove preemptible nodes from Google Kubernetes clusters before Google Cloud removes them the hard way.

## Problem solved

`gke-preemptible-sniper` helps in breaking down potentially big disruptions into smaller, more manageable ones. Instead of having a big chunk of your cluster removed at once, you can remove preemptible nodes one by one, giving your cluster time to recover and redistribute the load. This way, you can avoid the situation where your cluster is left with not enough resources to handle the load, since Google Clouds' preemption mechanism is not aware of the state of your cluster and does not necessarily respect disruption budgets of yours.

