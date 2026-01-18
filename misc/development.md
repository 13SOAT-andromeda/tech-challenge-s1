# ⚙️ Desenvolvimento local com kind

Este guia explica como instalar o kind, criar um cluster Kubernetes local e usar os comandos `make up` e `make down` disponíveis no `Makefile` do projeto para levantar e derrubar a infraestrutura local.

> Observação: os comandos e instruções abaixo foram adaptados ao ambiente do projeto (veja `Makefile`).

## Pré-requisitos

- Docker Desktop (ou Docker Engine) instalado e em execução. No WSL, use Docker Desktop com integração WSL.
- Make (GNU Make) instalado.
- Acesso ao terminal (WSL recomendado no Windows).

Ferramentas necessárias

- kind (Kubernetes IN Docker)
- kubectl

## Instalar kind

Ubuntu / WSL (bash):

```bash
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
```

Windows (PowerShell com Chocolatey):

```powershell
choco install kind
```

Ou via scoop:

```powershell
scoop install kind
```

Verifique a instalação:

```bash
kind --version
```

## Instalar kubectl

Ubuntu / WSL:

```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/
```

Verifique:

```bash
kubectl version --client
```

## Como funciona (Makefile)

O `Makefile` do projeto define os seguintes alvos relevantes:

- `make up` — cria um cluster kind, instala dependências (Ingress, Metrics Server), constrói a imagem Docker do app, carrega a imagem no cluster e aplica os manifests em `k8s/` (config, postgres, deployment, service, ingress, hpa). Também aguarda os pods ficarem prontos.
- `make down` — remove o cluster kind criado (executa `kind delete cluster --name tech-challenge-api-local`).

Os alvos dependentes são: `cluster`, `deps`, `build`, `load`, `deploy` (veja `Makefile` para detalhes). O `Makefile` usa `k8s/kind-config.yaml` para criar o cluster.

## Passos rápidos

1. Certifique-se que o Docker está em execução.
2. Execute:

```bash
make up
```

Espere até o Makefile terminar e a mensagem indicar que a aplicação está disponível.

Para derrubar a infra local criada pelo kind:

```bash
make down
```

## Verificações e depuração

- Verifique nodes e pods:

```bash
kubectl get nodes
kubectl get pods -A
```

- Verifique logs de um pod:

```bash
kubectl logs -n default <nome-do-pod>
```

- Se a imagem não for carregada, confira se o Docker daemon do host está acessível ao `kind`.

## Notas úteis

- O `Makefile` utiliza `kind create cluster --name tech-challenge-api-local --config k8s/kind-config.yaml`.
- Se já existir um cluster com o mesmo nome, execute `kind delete cluster --name tech-challenge-api-local` antes de criar novamente.
- O `make up` aplica o manifesto `k8s/postgres.yaml` para subir um Postgres dentro do cluster. Se preferir, você pode rodar seu próprio Postgres local e ajustar variáveis de ambiente da aplicação.

## Exemplos de comandos de verificação

```bash
# Verificar se o Ingress está pronto
kubectl get pods -n ingress-nginx

# Ver todos os serviços e ingress
kubectl get svc -A
kubectl get ingress -A
```

## Troubleshooting rápido

- Se o `make up` travar aguardando pods, rode `kubectl describe pod <pod> -n <namespace>` para identificar eventos.
- Se o Ingress não encaminhar corretamente, verifique o controller em `ingress-nginx` e as anotações do `Ingress` em `k8s/ingress.yaml`.
- Se houver problemas na criação do cluster, verifique a compatibilidade entre `kind`, `kubectl` e Docker.

---

Arquivo relacionado no repositório: `Makefile` (alvos `up` / `down`), `k8s/kind-config.yaml`, e os manifests em `k8s/`.

