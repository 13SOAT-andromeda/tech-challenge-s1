CLUSTER_NAME=tech-challenge-api-local
IMAGE_NAME=tech-challenge-api:latest

.PHONY: all up down deploy build load

up: cluster deps build load deploy
	@echo "Waiting for pods to be ready..."
	kubectl wait --for=condition=ready pod -l app=tech-challenge-api --timeout=60s
	@echo "✅ App is running at http://localhost/"

down:
	kind delete cluster --name $(CLUSTER_NAME)

cluster:
	kind create cluster --name $(CLUSTER_NAME) --config k8s/kind-config.yaml

deps:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

	@echo "Waiting Pods creating..."
	sleep 10

	@echo "Waiting Ingress Controller..."
	kubectl wait --namespace ingress-nginx \
	  --for=condition=ready pod \
	  --selector=app.kubernetes.io/component=controller \
	  --timeout=90s

	# 2. Install Metrics Server
	kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
	# Patch Metrics Server to work with Kind (Insecure TLS)
	kubectl patch -n kube-system deployment metrics-server --type=json \
	  -p '[{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--kubelet-insecure-tls"}]'

build:
	docker build -t $(IMAGE_NAME) .

load:
	kind load docker-image $(IMAGE_NAME) --name $(CLUSTER_NAME)

deploy:
	kubectl apply -f k8s/config.yaml
	kubectl apply -f k8s/postgres.yaml
	@echo "Waiting postgres starting..."
	kubectl wait --for=condition=ready pod -l app=postgres --timeout=90s
	kubectl apply -f k8s/deployment.yaml
	kubectl apply -f k8s/service.yaml
	kubectl apply -f k8s/ingress.yaml
	kubectl apply -f k8s/hpa.yaml