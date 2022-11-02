all: up build push apply

up: kubeup

dockerup:
	docker compose up -d

kubeup:
	./cluster_and_registry.sh

down: clusterdown dockerdown

clusterdown:
	kind delete cluster
	docker rm kind-registry -f

dockerdown:
	docker compose down

build:
	$(eval GIT_BRANCH=$(shell git rev-parse --short HEAD))
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-X main.ver=$(GIT_BRANCH)" -o testio .
	docker build -t localhost:5001/testio:latest .

push:
	docker push localhost:5001/testio:latest

apply:
	kubectl apply -f .k8s/specs/testio-deployment.yaml
	# kubectl apply -f .k8s/specs/testio-gateway.yaml
	# kubectl apply -f .k8s/specs/dashboard.yaml
	# kubectl apply -f .k8s/specs/service-account.yaml

delapps:
	kubectl delete -f .k8s/specs/testio-deployment.yaml

remove: delapps
	kubectl delete -f .k8s/specs/testio-gateway.yaml
	kubectl delete -f .k8s/specs/dashboard.yaml
	kubectl delete -f .k8s/specs/service-account.yaml

addons:
	kubectl apply -f .k8s/istio_addons/
	kubectl rollout status deployment/kiali -n istio-system

token:
	kubectl -n kubernetes-dashboard create token admin-user

migrate-up:
	migrate -path sql/migrations -database clickhouse://localhost:9000 up

migrate-down:
	migrate -path sql/migrations -database clickhouse://localhost:9000 down

kiali:
	istioctl dashboard kiali

linkerd:
	linkerd install --crds | kubectl apply -f -
	linkerd install | kubectl apply -f -
	linkerd check
	linkerd viz install | kubectl apply -f -
	linkerd check
	kubectl get deploy/ticktock-v1 -o yaml | linkerd inject - | kubectl apply -f -

grafana:
	helm repo add grafana https://grafana.github.io/helm-charts
	helm install grafana -n grafana --create-namespace grafana/grafana -f https://raw.githubusercontent.com/linkerd/linkerd2/main/grafana/values.yaml
	linkerd viz install --set grafana.url=grafana.grafana:3000 | kubectl apply -f -

ldash:
	linkerd viz dashboard