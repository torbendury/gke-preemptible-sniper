.RECIPEPREFIX = >
.PHONY: test unittest gohealth allchecks staticcheck govuln build helm

### Variables
RELEASE_IMAGE_NAME := torbendury/gke-preemptible-sniper

allchecks: test staticcheck govuln 

### Run the tests
test: unittest gohealth
> helm lint helm/gke-preemptible-sniper

unittest:
> go test -v -race ./...

gohealth:
> go mod verify
> go vet ./...

staticcheck:
> go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...

govuln:
> go run golang.org/x/vuln/cmd/govulncheck@latest ./...

### Build the release container
build:
> docker build --no-cache -t $(RELEASE_IMAGE_NAME):latest -f Dockerfile .

### Install the Helm Chart
helm:
> helm repo add gke-preemptible-sniper https://torbendury.github.io/gke-preemptible-sniper/ || true
> helm repo update
> helm upgrade --install gke-preemptible-sniper gke-preemptible-sniper/gke-preemptible-sniper --namespace gke-preemptible-sniper --create-namespace
