# gke-preemptible-sniper

`gke-preemptible-sniper` is a small application written in Golang that is supposed to run inside a Google Kubernetes cluster and work around the [known limitation of preemptible VMs](https://cloud.google.com/compute/docs/instances/preemptible#limitations).

Its' purpose is to gracefully remove preemptible nodes from Google Kubernetes clusters before Google Cloud removes them the hard way.

- [gke-preemptible-sniper](#gke-preemptible-sniper)
  - [Problem solved](#problem-solved)
  - [Installation](#installation)
    - [Helm](#helm)
  - [Metrics](#metrics)
  - [Development Status](#development-status)
  - [Resource Consumption](#resource-consumption)
  - [Testing](#testing)
  - [Example](#example)
    - [Infrastructure](#infrastructure)
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

Create a `values.yaml` file with the following content:

```yaml
serviceAccount:
  annotations:
    iam.gke.io/gcp-service-account: <SERVICE_ACCOUNT_NAME>@<PROJECT_ID>.iam.gserviceaccount.com
```

Install the chart:

```bash
helm install gke-preemptible-sniper gke-preemptible-sniper/gke-preemptible-sniper --namespace gke-preemptible-sniper --create-namespace --values=values.yaml
```

## Metrics

`gke-preemptible-sniper` provides Prometheus metrics on the `/metrics` endpoint. You can scrape them by configuring a Prometheus instance to scrape the metrics.

| Metric                                             | Description                                            |
|----------------------------------------------------|--------------------------------------------------------|
| `gke_preemptible_sniper_sniped_last_hour`          | Number of nodes sniped in the last hour                |
| `gke_preemptible_sniper_snipes_expected_next_hour` | Number of nodes expected to be sniped in the next hour |

Also, if you use Google Managed Prometheus or Prometheus Operator, you can configure the Helm Chart to automatically provide monitoring instrumentation for you. You can do this by adding the following to your `values.yaml`:

```yaml
metricScraping:
  googleManagedPrometheus: true
  # OR
  prometheusOperator: true
```

This will create a `ServiceMonitor` for Prometheus Operator or a `PodMonitoring` for Google Managed Prometheus.

## Development Status

This project is under active development. While I am using it in production, I cannot guarantee that it will work for you. If you encounter any issues, please open an issue on GitHub.

## Resource Consumption

`gke-preemptible-sniper` is designed to be lightweight and not consume too many resources.

| CPU Usage | Memory Usage | Container Image Size      |
|-----------|--------------|---------------------------|
| 0.001     | 10Mi         | 15MB (uncompressed: 55MB) |

## Testing

There are unit tests for the most important parts of the application.

Also, I e2e-tested the application by running it in a Google Kubernetes cluster and let it delete several preemptible nodes. Due to cost reasons this is not going to be part of the CI pipeline.

You can run the tests by executing the following command:

```bash
make test
```

This runs the unit tests, verifies go modules, `go vet` and also lints the Helm Chart.

## Example

### Infrastructure

You can find a working example of the infrastructure in the [terraform](terraform/) directory. It creates a Google Kubernetes cluster with everything needed around it. You can use it to test the application in a real environment, it might also serve as a starting point for your own infrastructure. I tested `gke-preemptible-sniper` with this infrastructure.

## Roadmap

For already released features, see the [changelog](CHANGELOG.md)! The following features are planned for future releases:

- [ ] gke-preemptible-sniper 1.1.0:
  - [ ] allow running outside cluster
  - [ ] read prepared kubeconfig

- [ ] gke-preemptible-sniper 1.2.0:
  - [ ] allow filtering out nodes by node label
  - [ ] stabilization: SIGTERM handling
