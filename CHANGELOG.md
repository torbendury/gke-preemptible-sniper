<a name="unreleased"></a>
## [Unreleased]

### Fix
- **google:** leave client closing to caller


<a name="gke-preemptible-sniper-0.2.8"></a>
## [gke-preemptible-sniper-0.2.8] - 2024-09-23
### Ci
- **Makefile:** speed up CI, leave out staticcheck/govuln on GH Action

### Fix
- **k8s:** ensure node drain timeout is greater than zero


<a name="gke-preemptible-sniper-0.2.7"></a>
## [gke-preemptible-sniper-0.2.7] - 2024-09-23
### Chore
- **main:** readable time until deletion

### Doc
- CHANGELOG
- **helm-chart:** version clarification


<a name="gke-preemptible-sniper-0.2.6"></a>
## [gke-preemptible-sniper-0.2.6] - 2024-09-23
### Fix
- **k8s-client:** don't throttle client loop


<a name="gke-preemptible-sniper-0.2.5"></a>
## [gke-preemptible-sniper-0.2.5] - 2024-09-23
### Fix
- **node-drain:** timeout settings similar to checkinterval


<a name="gke-preemptible-sniper-0.2.4"></a>
## [gke-preemptible-sniper-0.2.4] - 2024-09-23
### Chore
- **helm-chart:** install notes

### Fix
- **checkinterval:** ensure interval > 0


<a name="gke-preemptible-sniper-0.2.3"></a>
## [gke-preemptible-sniper-0.2.3] - 2024-09-23
### Fix
- **k8s:** untighten client rate limiting


<a name="gke-preemptible-sniper-0.2.2"></a>
## [gke-preemptible-sniper-0.2.2] - 2024-09-23
### Fix
- **main:** context cancellations


<a name="gke-preemptible-sniper-0.2.1"></a>
## [gke-preemptible-sniper-0.2.1] - 2024-09-22
### Ci
- initial dependabot config
- **gh-actions:** run unit tests on dependabot PRs
- **gh-actions:** move from state outputs to new patched version

### Doc
- CHANGELOG
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


[Unreleased]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.2.8...HEAD
[gke-preemptible-sniper-0.2.8]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.2.7...gke-preemptible-sniper-0.2.8
[gke-preemptible-sniper-0.2.7]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.2.6...gke-preemptible-sniper-0.2.7
[gke-preemptible-sniper-0.2.6]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.2.5...gke-preemptible-sniper-0.2.6
[gke-preemptible-sniper-0.2.5]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.2.4...gke-preemptible-sniper-0.2.5
[gke-preemptible-sniper-0.2.4]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.2.3...gke-preemptible-sniper-0.2.4
[gke-preemptible-sniper-0.2.3]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.2.2...gke-preemptible-sniper-0.2.3
[gke-preemptible-sniper-0.2.2]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.2.1...gke-preemptible-sniper-0.2.2
[gke-preemptible-sniper-0.2.1]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.2.0...gke-preemptible-sniper-0.2.1
[gke-preemptible-sniper-0.2.0]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.1.3...gke-preemptible-sniper-0.2.0
[gke-preemptible-sniper-0.1.3]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.1.1...gke-preemptible-sniper-0.1.3
[gke-preemptible-sniper-0.1.1]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.1.0...gke-preemptible-sniper-0.1.1
[gke-preemptible-sniper-0.1.0]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.0.0...gke-preemptible-sniper-0.1.0
[gke-preemptible-sniper-0.0.0]: https://github.com/torbendury/kube-networkpolicy-denier/compare/0.0.0...gke-preemptible-sniper-0.0.0
