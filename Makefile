include .env

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

push:
	@echo "Pushing docker image into AWS..."
	@echo "Docker and AWS CLI are needed before push API image..."
	docker build --target production -t $(AWS_ACCOUNT).dkr.ecr.$(AWS_REGION).amazonaws.com/$(AWS_ECR_REPO):latest .
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin ${AWS_ACCOUNT}.dkr.ecr.${AWS_REGION}.amazonaws.com
	docker push $(AWS_ACCOUNT).dkr.ecr.$(AWS_REGION).amazonaws.com/$(AWS_ECR_REPO):latest

load:
	kind load docker-image $(IMAGE_NAME) --name $(CLUSTER_NAME)

deploy-local:
	@echo "Deploying local..."
	kubectl apply -k k8s/overlays/local
	@echo "Aguardando a inicialização do posgres..."
	kubectl wait --for=condition=ready pod -l app=postgres --timeout=90s || true

deploy-aws:
	@RDS_ADDRESS=$$(terraform output -raw rds_address 2>/dev/null) && \
	if [ -z "$$RDS_ADDRESS" ]; then \
		echo "Erro: Não foi possível obter o endereço do RDS. Verifique se o Terraform foi aplicado."; \
		exit 1; \
	fi
	@echo "Atualizando o endereço do banco de dados para o endpoint do RDS"
	kubectl patch configmap api-config \
		--type merge \
		-p "{\"data\":{\"DB_HOST\":\"$$RDS_ADDRESS\"}}" && \
	@echo "Aplicando a configuração do cluster"
	kubectl apply -k k8s/overlays/aws
	@echo "IP Externo..."
	kubectl get svc tech-challenge-api-svc | awk '{print $$4}'

deploy:
ifeq ($(ENV),aws)
	$(MAKE) deploy-aws
else
	$(MAKE) deploy-local
endif