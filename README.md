# gke-preemptible-sniper

`gke-preemptible-sniper` is a small application written in Golang that is supposed to run inside a Google Kubernetes cluster and work around the [known limitation of preemptible VMs](https://cloud.google.com/compute/docs/instances/preemptible#limitations).

Its' purpose is to gracefully remove preemptible nodes from Google Kubernetes clusters before Google Cloud removes them the hard way.

- [gke-preemptible-sniper](#gke-preemptible-sniper)
  - [Problem solved](#problem-solved)
  - [Installation](#installation)
    - [Helm](#helm)
  - [Metrics](#metrics)
  - [Development Status](#development-status)
  - [Testing](#testing)
  - [Roadmap](#roadmap)

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

## Metrics

`gke-preemptible-sniper` provides Prometheus metrics on the `/metrics` endpoint. You can scrape them by configuring a Prometheus instance to scrape the metrics.

| Metric                                             | Description                                            |
|----------------------------------------------------|--------------------------------------------------------|
| `gke_preemptible_sniper_sniped_last_hour`          | Number of nodes sniped in the last hour                |
| `gke_preemptible_sniper_snipes_expected_next_hour` | Number of nodes expected to be sniped in the next hour |

## Development Status

This project is in a very early stage of development. It is not recommended to use it in production environments yet.

## Testing

There are unit tests for the most important parts of the application.

Also, I e2e-tested the application by running it in a Google Kubernetes cluster and let it delete several preemptible nodes. Due to cost reasons this is not going to be part of the CI pipeline.

## Roadmap

For already released features, see the [changelog](CHANGELOG.md)! The following features are planned for future releases:

- [ ] gke-preemptible-sniper 1.0.0:
  - [ ] stabilization
  - [x] sensible defaults
  - [ ] documentation
  - [ ] logging improvements
  - [ ] error budget improvements

- [ ] gke-preemptible-sniper 1.1.0:
  - [ ] allow running outside cluster
  - [ ] read prepared kubeconfig

- [ ] gke-preemptible-sniper 1.2.0:
  - [ ] allow filtering out nodes by node label
