<a name="unreleased"></a>
## [Unreleased]

### Chore
- typo

### Doc
- title picture


<a name="gke-preemptible-sniper-1.1.3"></a>
## [gke-preemptible-sniper-1.1.3] - 2025-01-16
### Chore
- bump version
- bump dependencies


<a name="gke-preemptible-sniper-1.1.2"></a>
## [gke-preemptible-sniper-1.1.2] - 2024-11-04
### Fix
- **stats:** reset metric before adding new records


<a name="gke-preemptible-sniper-1.1.1"></a>
## [gke-preemptible-sniper-1.1.1] - 2024-10-31
### Fix
- **stats:** sniped nodes stats calculation


<a name="gke-preemptible-sniper-1.1.0"></a>
## [gke-preemptible-sniper-1.1.0] - 2024-10-29
### Chore
- bump version

### Feat
- golang 1.23 and deps


<a name="gke-preemptible-sniper-1.0.1"></a>
## [gke-preemptible-sniper-1.0.1] - 2024-10-28
### Chore
- bump version

### Doc
- sniping behavior
- CHANGELOG and project docs

### Fix
- **k8s:** do not leave nodes dangling if eviction does not work


<a name="gke-preemptible-sniper-1.0.0"></a>
## [gke-preemptible-sniper-1.0.0] - 2024-09-24
### Chore
- bump version

### Doc
- resource consumption
- development status


<a name="gke-preemptible-sniper-0.4.7"></a>
## [gke-preemptible-sniper-0.4.7] - 2024-09-24
### Feat
- **monitoring:** add PodMonitoring and ServiceMonitor for auto instrumented metric scraping


<a name="gke-preemptible-sniper-0.4.6"></a>
## [gke-preemptible-sniper-0.4.6] - 2024-09-24
### Chore
- **error budget:** re-check at end of loop and increase initial error budget

### Fix
- **gcloud:** implement up to 3 retries for read-only methods


<a name="gke-preemptible-sniper-0.4.5"></a>
## [gke-preemptible-sniper-0.4.5] - 2024-09-24
### Chore
- typo

### Doc
- document go packages and their functions
- roadmap status
- main loop documentation
- **helm-chart:** description update


<a name="gke-preemptible-sniper-0.4.4"></a>
## [gke-preemptible-sniper-0.4.4] - 2024-09-24
### Refactor
- get rid of magic numbers inside code


<a name="gke-preemptible-sniper-0.4.3"></a>
## [gke-preemptible-sniper-0.4.3] - 2024-09-24
### Chore
- encapsulate error budget handling and error handling/logging


<a name="gke-preemptible-sniper-0.4.2"></a>
## [gke-preemptible-sniper-0.4.2] - 2024-09-24
### Fix
- **main loop:** wait for concurrent node processing to finish before context cancellation


<a name="gke-preemptible-sniper-0.4.1"></a>
## [gke-preemptible-sniper-0.4.1] - 2024-09-24
### Fix
- **main loop:** give each goroutine its own context


<a name="gke-preemptible-sniper-0.4.0"></a>
## [gke-preemptible-sniper-0.4.0] - 2024-09-23
### Doc
- contributing
- CODE_OF_CONDUCT
- CHANGELOG
- roadmap scope

### Feat
- **main:** process nodes concurrently


<a name="gke-preemptible-sniper-0.3.5"></a>
## [gke-preemptible-sniper-0.3.5] - 2024-09-23
### Fix
- **helm-chart:** auto-formatter breaking helm templates
- **rbac:** permit deletion of kubernetes nodes


<a name="gke-preemptible-sniper-0.3.4"></a>
## [gke-preemptible-sniper-0.3.4] - 2024-09-23
### Chore
- less verbose logging


<a name="gke-preemptible-sniper-0.3.3"></a>
## [gke-preemptible-sniper-0.3.3] - 2024-09-23
### Feat
- **k8s:** allow k8s node deletion together with GCP instance deletion


<a name="gke-preemptible-sniper-0.3.2"></a>
## [gke-preemptible-sniper-0.3.2] - 2024-09-23
### Feat
- **logging:** log consolidation and improvements


<a name="gke-preemptible-sniper-0.3.1"></a>
## [gke-preemptible-sniper-0.3.1] - 2024-09-23
### Debug
- **metrics:** metrics update logging

### Doc
- 1.0.0 roadmap


<a name="gke-preemptible-sniper-0.3.0"></a>
## [gke-preemptible-sniper-0.3.0] - 2024-09-23
### Chore
- gomod update
- log successfully deleted instances + sleep mode

### Doc
- TOC
- roadmap scope
- roadmap scopes
- metrics
- testing
- CHANGELOG

### Feat
- **main:** let errors be handled by error budget functionality
- **main:** introduce error budget errors cause the budget to decrease over time, after a decrease the application will try to heal itself
- **metrics:** enable scrapable prometheus metrics

### Refactor
- **main:** separate out checking error budget in main loop
- **main:** separate out loop functionality


<a name="gke-preemptible-sniper-0.2.12"></a>
## [gke-preemptible-sniper-0.2.12] - 2024-09-23
### Fix
- **helm-chart:** workload identity annotation


<a name="gke-preemptible-sniper-0.2.11"></a>
## [gke-preemptible-sniper-0.2.11] - 2024-09-23
### Chore
- bump version

### Feat
- **helm-chart:** use workload identity


<a name="gke-preemptible-sniper-0.2.10"></a>
## [gke-preemptible-sniper-0.2.10] - 2024-09-23
### Debug
- **main:** check for empty instance and zone labels

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


[Unreleased]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-1.1.3...HEAD
[gke-preemptible-sniper-1.1.3]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-1.1.2...gke-preemptible-sniper-1.1.3
[gke-preemptible-sniper-1.1.2]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-1.1.1...gke-preemptible-sniper-1.1.2
[gke-preemptible-sniper-1.1.1]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-1.1.0...gke-preemptible-sniper-1.1.1
[gke-preemptible-sniper-1.1.0]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-1.0.1...gke-preemptible-sniper-1.1.0
[gke-preemptible-sniper-1.0.1]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-1.0.0...gke-preemptible-sniper-1.0.1
[gke-preemptible-sniper-1.0.0]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.4.7...gke-preemptible-sniper-1.0.0
[gke-preemptible-sniper-0.4.7]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.4.6...gke-preemptible-sniper-0.4.7
[gke-preemptible-sniper-0.4.6]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.4.5...gke-preemptible-sniper-0.4.6
[gke-preemptible-sniper-0.4.5]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.4.4...gke-preemptible-sniper-0.4.5
[gke-preemptible-sniper-0.4.4]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.4.3...gke-preemptible-sniper-0.4.4
[gke-preemptible-sniper-0.4.3]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.4.2...gke-preemptible-sniper-0.4.3
[gke-preemptible-sniper-0.4.2]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.4.1...gke-preemptible-sniper-0.4.2
[gke-preemptible-sniper-0.4.1]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.4.0...gke-preemptible-sniper-0.4.1
[gke-preemptible-sniper-0.4.0]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.3.5...gke-preemptible-sniper-0.4.0
[gke-preemptible-sniper-0.3.5]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.3.4...gke-preemptible-sniper-0.3.5
[gke-preemptible-sniper-0.3.4]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.3.3...gke-preemptible-sniper-0.3.4
[gke-preemptible-sniper-0.3.3]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.3.2...gke-preemptible-sniper-0.3.3
[gke-preemptible-sniper-0.3.2]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.3.1...gke-preemptible-sniper-0.3.2
[gke-preemptible-sniper-0.3.1]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.3.0...gke-preemptible-sniper-0.3.1
[gke-preemptible-sniper-0.3.0]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.2.12...gke-preemptible-sniper-0.3.0
[gke-preemptible-sniper-0.2.12]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.2.11...gke-preemptible-sniper-0.2.12
[gke-preemptible-sniper-0.2.11]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.2.10...gke-preemptible-sniper-0.2.11
[gke-preemptible-sniper-0.2.10]: https://github.com/torbendury/kube-networkpolicy-denier/compare/gke-preemptible-sniper-0.2.8...gke-preemptible-sniper-0.2.10
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
