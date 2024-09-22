<a name="unreleased"></a>
## [Unreleased]

### Doc
- CHANGELOG
- **readme:** new roadmap item

### Hotfix
- bump version
- **helm:** ensure interval and timeout seconds to be strings


<a name="gke-preemptible-sniper-0.2.0"></a>
## [gke-preemptible-sniper-0.2.0] - 2024-09-22
### Chore
- ignore test coverage files

### Doc
- allowlist/blocklist hours
- **helm-chart:** installation note command

### Feat
- **core:** configurable node drain timeout
- **core:** greater timeout for actual draining nodes
- **time:** configurable main loop check interval
- **time:** add optional feature for allowed and blocked hours


<a name="gke-preemptible-sniper-0.1.3"></a>
## [gke-preemptible-sniper-0.1.3] - 2024-09-22
### Chore
- bump version

### Feat
- **core:** health check - failures to be implemented
- **k8s:** check for label instead of annotation


<a name="gke-preemptible-sniper-0.1.1"></a>
## [gke-preemptible-sniper-0.1.1] - 2024-09-22
### Chore
- bump version

### Ci
- **Makefile:** option to install helm chart

### Feat
- **infra:** test infra for helm chart
- **terraform:** make workload identity stuff depend on GKE which is provisioning the WI pool

### Fix
- **main:** context cancellation


<a name="gke-preemptible-sniper-0.1.0"></a>
## [gke-preemptible-sniper-0.1.0] - 2024-09-22
### Chore
- bump version

### Chroe
- **gcloud:** remove non-mockable client creation

### Doc
- **helm-chart:** document installation path

### Feat
- **gcloud:** rewrite to use new library
- **main loop:** don't exit on minor errors, retry later. also use JSON logging


<a name="gke-preemptible-sniper-0.0.0"></a>
## [gke-preemptible-sniper-0.0.0] - 2024-09-22
### Chore
- ignore workspace settings

### Ci
- naming and basic pipelines, make them work later
- initial Dockerfile
- **Makefile:** sanitize Makefile for a good default
- **gh-actions:** refine naming
- **helm-chart:** initial structure, to be refined later
- **pre-commit:** initial config with changelogging

### Doc
- formatting
- formatting
- roadmap and project status

### Feat
- **helm-template:** correct RBAC permissions

### Fix
- **k8s_test:** do not expect kubeconfig fallback to work in CI
- **k8s_test:** dont expect kubeconfig fallback to work in CI

### Test
- **k8s-client:** unit tests for all functions
- **k8s-client:** mock k8s client


<a name="0.0.0"></a>
## 0.0.0 - 2024-09-16
### Ci
- initial CI files, to be adjusted later

### Feat
- first working version for 0.0.0
- initial code structure

### Refactor
- move projectId retrieval out
- different clients for kubernetes and GCP, move init into their own function


[Unreleased]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.2.0...HEAD
[gke-preemptible-sniper-0.2.0]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.1.3...gke-preemptible-sniper-0.2.0
[gke-preemptible-sniper-0.1.3]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.1.1...gke-preemptible-sniper-0.1.3
[gke-preemptible-sniper-0.1.1]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.1.0...gke-preemptible-sniper-0.1.1
[gke-preemptible-sniper-0.1.0]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.0.0...gke-preemptible-sniper-0.1.0
[gke-preemptible-sniper-0.0.0]: https://github.com/torbendury/kube-networkpolicy-denier/compare/0.0.0...gke-preemptible-sniper-0.0.0
