# ⚙️ Tech Challenge S1 - API de Gestão de Oficina Mecânica

Este projeto é uma API completa para gerenciamento de uma oficina mecânica, desenvolvida em Go (Golang). Ele permite gerenciar clientes, veículos, produtos e ordens de serviço, com autenticação JWT, documentação via Swagger e implantação simplificada com Docker Compose.

## 🚀 Como Executar

### Com Docker (Recomendado)

O repositório inclui `Dockerfile` e `docker-compose.yml` para iniciar a aplicação e o banco de dados automaticamente.

1.  **Instale o Docker**:
    - Windows: [Docker Desktop](https://docs.docker.com/desktop/)
    - Ubuntu:
      ```bash
      sudo apt update && sudo apt install -y docker.io docker-compose
      sudo systemctl enable --now docker
      ```
2.  Crie um arquivo `.env` na raiz do projeto. Para desenvolvimento, você pode usar os seguintes valores:
    ```env
    DB_HOST=db
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=garagedb
    DB_PORT=5432
    DB_SSLMODE=disable
    DB_TIMEZONE=UTC
    ENV=development
    HTTP_ALLOWED_ORIGINS=*
    HTTP_PORT=8080
    ADMIN_EMAIL=admin@admin.com
    ADMIN_PASSWORD=Admin123!
    MAILTRAP_TOKEN=6a45f171cfc233e4edc93d8b847cf19f
    MAILTRAP_URL="https://send.api.mailtrap.io/api"
    ```
3.  Inicie os serviços:
    ```bash
    docker-compose up --build
    ```
    Ou em segundo plano:
    ```bash
    docker-compose up --build -d
    ```
4.  Para parar os serviços:
    ```bash
    docker-compose down
    ```

### 🏗️ Infraestrutura

Esta seção descreve as opções de infraestrutura que você pode usar localmente para desenvolver e testar a aplicação.

Opções principais:

- Docker Compose (rápido): já configurado no repositório. Executando `docker-compose up --build` você terá a API e um banco (se configurado no compose) prontos para uso.

- Kubernetes local com `kind` (mais próximo do ambiente de produção): o projeto inclui manifests em `k8s/` e um `Makefile` com alvos para criar um cluster kind, construir e carregar a imagem e aplicar os manifests.

Comandos úteis (na raiz do projeto):

```bash
# Subir infra rápida com Docker Compose
docker-compose up --build

# Criar cluster local, construir imagem, carregar e aplicar manifests
make up

# Derrubar o cluster criado pelo Makefile
make down
```

Onde procurar os manifests e a configuração do cluster:

- Manifests Kubernetes: `k8s/`
- Arquivo de configuração do kind usado pelo Makefile: `k8s/kind-config.yaml`
- Guia passo-a-passo para criar o cluster local (instalação do kind, kubectl, notas): [development.md](./development.md)

Notas rápidas:

- `make up` depende de ter `kind`, `kubectl` e `docker` instalados e em funcionamento (veja [development.md](./development.md) para detalhes de instalação).
- Se preferir um Postgres local em contêiner, use o comando de exemplo na seção "Banco de dados rápido com Docker".

### Localmente (sem Docker)

**Pré-requisitos**:
- Go 1.25+ instalado.
- Uma instância de PostgreSQL em execução (local ou remota) acessível com as credenciais fornecidas.

**Passos:**

1.  **Instale o Go**:
    -   **Windows**: Baixe e execute o instalador em [https://go.dev/dl/](https://go.dev/dl/).
    -   **Ubuntu (WSL)**:
        ```bash
        sudo apt update && sudo apt install -y golang
        ```

2.  **Crie um arquivo `.env` na raiz do projeto**:
    ```env
    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=postgres123
    DB_NAME=myapp_db
    DB_PORT=5432
    DB_SSLMODE=disable
    DB_TIMEZONE=America/Sao_Paulo
    ENV=production
    HTTP_ALLOWED_ORIGINS=*
    GIN_MODE=release
    JWT_SECRET=5b9b178c235820c6e69fbf54876bc4df3ffb4f3ab5ec87305b8b42d2481358c3
    JWT_ACCESS_TOKEN_EXPIRY=24h
    JWT_REFRESH_TOKEN_EXPIRY=700h
    ADMIN_PASSWORD=Admin123!
    ADMIN_EMAIL=admin@admin.com
    MAILTRAP_TOKEN=6a45f171cfc233e4edc93d8b847cf19f
    MAILTRAP_URL="https://send.api.mailtrap.io/api"
    ```
    *Nota: As variáveis são carregadas em `internal/adapter/config/config.go`.*

3.  **Execute a aplicação**:
    Pode-se executar diretamente com `go run`:
    ```bash
    go run ./cmd/api
    ```
    Ou compilar e depois executar:
    ```bash
    go build -o bin/app ./cmd/api
    ./bin/app
    ```

4.  **Verifique se a aplicação está em execução**:
    -   **Health Check**: [http://localhost:8080/health](http://localhost:8080/health)
    -   **Swagger UI**: [http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)
    -   **Redoc**: [http://localhost:8080/redoc](http://localhost:8080/redoc)

---
#### **Banco de dados rápido com Docker (Opcional)**
Caso não tenha o PostgreSQL instalado, você pode rodar um contêiner temporário. Lembre-se de ajustar o `DB_HOST` no seu `.env` para `localhost`.

```bash
docker run --name tech-challenge-pg \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres123 \
  -e POSTGRES_DB=myapp_db \
  -p 5432:5432 -d postgres:15
```

Para parar e remover o contêiner:
```bash
docker stop tech-challenge-pg && docker rm tech-challenge-pg
```
---

## 🧪 Como Rodar os Testes E2E

**Pré-requisitos**: A aplicação deve estar rodando localmente conforme a etapa: `Executar localmente (sem Docker)`.

Na raiz do projeto, executar o comando:
```bash
go test ./test/e2e/... -v
```

## 📚 Documentação da API

Você pode acessar a documentação interativa da API através dos seguintes endpoints:

-   **Swagger UI**: [http://localhost:8080/docs/index.html?url=/swagger/swagger.yaml](http://localhost:8080/docs/index.html?url=/swagger/swagger.yaml)
-   **Redoc UI**: [http://localhost:8080/redoc/](http://localhost:8080/redoc/)

A especificação OpenAPI está disponível em `/swagger/swagger.yaml`.

## 🏗️ Arquitetura do Projeto

O projeto adota a **Arquitetura Hexagonal** (Ports and Adapters) para isolar a lógica de negócios (o núcleo da aplicação) de detalhes de infraestrutura e frameworks externos. Isso garante que o núcleo permaneça puro e testável, enquanto as tecnologias externas podem ser trocadas sem impactar as regras de negócio.

O fluxo de dependência é sempre direcionado para o centro do hexágono.

```
       +------------------------------------------------+
       |                 Driving Adapters               |
       | (HTTP Handlers, CLI, Testes de Aceitação)      |
       +-----------------------+------------------------+
                               |
                               v (Driving Ports)
+-------------------------------------------------------------+
|                      Application Core                         |
|                                                             |
|  +---------------------+      +--------------------------+  |
|  | Application Services|----->|      Domain Objects      |  |
|  |   (Use Cases)       |      | (Entidades, Value Objects)|  |
|  +---------------------+      +--------------------------+  |
|                                                             |
+------------------------+------------------------------------+
                         |
                         v (Driven Ports)
       +-----------------------+------------------------+
       |                  Driven Adapters               |
       | (Database Repositories, Email Services, APIs)  |
       +------------------------------------------------+

```

### Componentes Principais no Projeto:

-   **🎯 **Core (Núcleo da Aplicação)**: O centro da aplicação, onde a lógica de negócio reside.
    -   **`internal/domain`**: Contém as entidades de negócio (`Customer`, `Order`, `Product`), objetos de valor (`vo_document`, `vo_plate`) e as regras de negócio mais puras. Esta camada não depende de nenhuma outra.
    -   **`internal/application/services`**: Implementa os casos de uso da aplicação (ex: `CompanyService`, `CustomerService`). Orquestra as entidades de domínio e utiliza as `ports` (interfaces) para interagir com o mundo exterior, sem conhecer os detalhes da implementação.

-   **🔌 **Ports (Portas)**: São as interfaces que definem os contratos de comunicação.
    -   **`internal/application/ports`**: Este diretório é crucial, pois define todas as `ports`.
        -   **Driving Ports (Portas Primárias)**: São as interfaces dos próprios serviços (ex: `CompanyService`, `CustomerService`), que definem como os adaptadores primários podem interagir com a aplicação.
        -   **Driven Ports (Portas Secundárias)**: São as interfaces que o núcleo da aplicação precisa para se comunicar com serviços externos, como `CompanyRepository`, `CustomerRepository` e `Email`.

-   **🔩 **Adapters (Adaptadores)**: Implementam as `ports` para conectar o núcleo com o mundo exterior.
    -   **`internal/adapter`**: Contém todas as implementações concretas.
        -   **Driving Adapters (Adaptadores Primários)**: Invocam a aplicação.
            -   **`adapter/http/handlers`**: Recebem requisições HTTP, validam os dados e chamam os serviços da aplicação (o núcleo).
            -   **`cmd/api/main.go`**: Ponto de entrada que inicializa e conecta todos os adaptadores e serviços (injeção de dependência).
        -   **Driven Adapters (Adaptadores Secundários)**: São "controlados" pela aplicação.
            -   **`adapter/database/repository`**: Implementam as interfaces de repositório (`Driven Ports`) para persistir dados no PostgreSQL.
            -   **`adapter/email`**: Implementa a interface `Email`, atuando como um cliente para a API de envio de e-mails do Mailtrap.

-   **📁 `pkg`**: Contém pacotes reutilizáveis e agnósticos ao negócio, como criptografia (`encryption`), manipulação de JWT (`jwt`) e conversores (`converters`).

## 🛠️ Definições Técnicas

Esta seção descreve as principais escolhas tecnológicas feitas para este projeto.

-   **🐘 Banco de Dados (PostgreSQL)**: Escolhido por sua robustez, confiabilidade e suporte a consultas complexas. É ideal para sistemas transacionais que exigem alta integridade de dados, como o gerenciamento de uma oficina.

-   **🐹 Linguagem (Golang)**: Selecionada por seu alto desempenho, simplicidade e excelente suporte à concorrência. O Go permite construir APIs rápidas, eficientes e escaláveis, com um baixo consumo de recursos.

-   **🏛️ Arquitetura (Hexagonal)**: Adotada para isolar a lógica de negócios das dependências externas. Isso aumenta a testabilidade, facilita a manutenção e permite trocar tecnologias futuras (como o banco de dados) sem impactar o núcleo da aplicação.

-   **🐳 Containerização (Docker)**: Utilizado para padronizar e simplificar o ambiente de desenvolvimento e implantação. Com o Docker, a aplicação e suas dependências são executadas de forma consistente em qualquer ambiente com um único comando.

## 🔧 Solução de problemas

Se a aplicação falhar ao conectar-se ao banco de dados, verifique se o seu Postgres está em execução e se os valores do `.env` correspondem. Verifique os logs impressos em `stdout`/`stderr` para erros de inicialização.

## 🔗 Misc

-   **Postman Collection**: Você pode fazer o download da collection do Postman para testar os endpoints da API [aqui](./misc/Tech%20Challenge%20S1.postman.json).
-   **Guia de Desenvolvimento (kind & cluster local)**: [development.md](./development.md)
