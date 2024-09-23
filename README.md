# gke-preemptible-sniper

`gke-preemptible-sniper` is a small application written in Golang that is supposed to run inside a Google Kubernetes cluster and work around the [known limitation of preemptible VMs](https://cloud.google.com/compute/docs/instances/preemptible#limitations).

Its' purpose is to gracefully remove preemptible nodes from Google Kubernetes clusters before Google Cloud removes them the hard way.

## Problem solved

`gke-preemptible-sniper` helps in breaking down potentially big disruptions into smaller, more manageable ones. Instead of having a big chunk of your cluster removed at once, you can remove preemptible nodes one by one, giving your cluster time to recover and redistribute the load. This way, you can avoid the situation where your cluster is left with not enough resources to handle the load, since Google Clouds' preemption mechanism is not aware of the state of your cluster and does not necessarily respect disruption budgets of yours.

## Installation

### Helm

Add the repository to your local Helm installation:

```bash
helm repo add gke-preemptible-sniper https://torbendury.github.io/gke-preemptible-sniper
helm repo update
```

Install the chart:

```bash
helm install gke-preemptible-sniper gke-preemptible-sniper/gke-preemptible-sniper --namespace gke-preemptible-sniper --create-namespace
```

## Development Status

This project is in a very early stage of development. It is not recommended to use it in production environments yet.

## Testing

There are unit tests for the most important parts of the application.

Also, I e2e-tested the application by running it in a Google Kubernetes cluster and let it delete several preemptible nodes. Due to cost reasons this is not going to be part of the CI pipeline.

## Roadmap

- [x] gke-preemptible-sniper 0.0.0:
  - [x] retrieve running nodes in cluster
  - [x] drain node
  - [x] cordon node
  - [x] retrieve zone and GCP project id
  - [x] check if node has annotation set
  - [x] set annotation
  - [x] get VM instance in GCP
  - [x] delete VM instance in GCP

- [x] gke-preemptible-sniper 0.1.0:
  - [x] Containerize
  - [x] Publish to DockerHub
  - [x] Helm Chart
  - [x] GitHub Actions Build

- [x] gke-preemptible-sniper 0.2.0:
  - [x] allowlist hours
  - [x] blocklist hours
  - [x] configurable check interval
  - [x] configurable node drain timeout

- [ ] gke-preemptible-sniper 0.3.0:
  - [ ] provide prometheus metrics
  - [ ] allow filtering out nodes by node label
  - [ ] allow filtering out nodes by pod label
  - [x] introduce error budget

- [ ] gke-preemptible-sniper 0.4.0:
  - allow running outside cluster
  - read kubeconfig

- [ ] gke-preemptible-sniper 0.5.0:
  - [ ] process nodes concurrently
