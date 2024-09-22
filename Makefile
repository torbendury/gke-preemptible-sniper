.RECIPEPREFIX = >
.PHONY: test unittest gohealth build helm

### Variables
RELEASE_IMAGE_NAME := torbendury/gke-preemptible-sniper

### Run the tests
test: unittest gohealth
> helm lint helm/gke-preemptible-sniper

unittest:
> go test -v -race ./...

gohealth:
> go mod verify
> go vet ./...
> go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
> go run golang.org/x/vuln/cmd/govulncheck@latest ./...

### Build the release container
build:
> docker build --no-cache -t $(RELEASE_IMAGE_NAME):latest -f Dockerfile .

### Install the Helm Chart
helm:
> helm upgrade --install gke-preemptible-sniper ./helm/gke-preemptible-sniper --set image.tag=latest --namespace gke-preemptible-sniper --create-namespace
