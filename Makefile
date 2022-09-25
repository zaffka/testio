.PHONY: cluster shut build

cluster:
	./cluster_and_registry.sh

shut:
	kind delete cluster

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-X main.periodStr=5s" -o testio .
	docker build -t localhost:5000/testio:latest .