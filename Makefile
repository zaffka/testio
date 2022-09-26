.PHONY: up down build apply remove

up:
	./cluster_and_registry.sh

down:
	kind delete cluster
	docker rm kind-registry -f

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-X main.periodStr=5s" -o testio .
	docker build -t localhost:5001/testio:latest .
	docker push localhost:5001/testio:latest

apply:
	kubectl apply -f testio.yaml

remove:
	kubectl delete -f testio.yaml