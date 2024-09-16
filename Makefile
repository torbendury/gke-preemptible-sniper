.RECIPEPREFIX = >
.PHONY: test build savebuild kube kuberun github helm local localproxy

### Variables
RELEASE_IMAGE_NAME := torbendury/gke-preemptible-sniper

### Run the tests
test:
> go test -v -race ./...
> helm lint helm/gke-preemptible-sniper
> go mod verify
> go vet ./...
> go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
> go run golang.org/x/vuln/cmd/govulncheck@latest ./...

### Build the release container
build:
> docker build --no-cache -t $(RELEASE_IMAGE_NAME):latest --target release .


savebuild:
> docker image save -o image.tar $(RELEASE_IMAGE_NAME):latest

### Create local test cluster
kube: savebuild
> minikube start
> minikube -p minikube docker-env
> minikube image load image.tar
> sleep 10

kuberun:
> kubectl run --rm -i gke-preemptible-sniper --image=$(RELEASE_IMAGE_NAME):latest --image-pull-policy=IfNotPresent

### Install the Helm Chart
helm:
> helm upgrade --install gke-preemptible-sniper ./helm/gke-preemptible-sniper --set image.tag=latest --namespace gke-preemptible-sniper --create-namespace

### Create a complete local environment
local: build kube helm
