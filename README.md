# gke-preemptible-sniper

`gke-preemptible-sniper` is a small application written in Golang that is supposed to run inside a Google Kubernetes cluster and work around the [known limitation of preemptible VMs](https://cloud.google.com/compute/docs/instances/preemptible#limitations).

Its' purpose is to gracefully remove preemptible nodes from Google Kubernetes clusters before Google Cloud removes them the hard way.

## Problem solved

`gke-preemptible-sniper` helps in breaking down potentially big disruptions into smaller, more manageable ones. Instead of having a big chunk of your cluster removed at once, you can remove preemptible nodes one by one, giving your cluster time to recover and redistribute the load. This way, you can avoid the situation where your cluster is left with not enough resources to handle the load, since Google Clouds' preemption mechanism is not aware of the state of your cluster and does not necessarily respect disruption budgets of yours.

## Development Status

This project is in a very early stage of development. It is not recommended to use it in production environments yet.

## Roadmap

- [x] gke-preemptible-sniper 0.0.0:
    ✔ retrieve running nodes in cluster @done(24-09-16 22:23)
    ✔ drain node @done(24-09-16 22:27)
    ✔ cordon node @done(24-09-16 22:25)
    ✔ retrieve zone and GCP project id @done(24-09-16 22:34)
    ✔ check if node has annotation set @done(24-09-16 22:28)
    ✔ set annotation @done(24-09-16 22:28)
    ✔ get VM instance in GCP @done(24-09-16 22:33)
    ✔ delete VM instance in GCP @done(24-09-16 22:52)


- [ ] gke-preemptible-sniper 0.1.0:
  - Containerize
  - Publish to DockerHub
  - Helm Chart
  - GitHub Actions Build

- [ ] gke-preemptible-sniper 0.2.0:
  - allowlist hours
  - blocklist hours
  - configurable check interval
  - configurable node drain timeout

- [ ] gke-preemptible-sniper 0.3.0:
  - provide prometheus metrics
  - allow filtering out nodes by node label
  - allow filtering out nodes by pod label

- [ ] gke-preemptible-sniper 0.4.0:
  - allow running outside cluster
  - read kubeconfig
