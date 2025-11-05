# Tech Challenge S1

Uma API completa para gestão de oficina mecânica, desenvolvida em Go (Golang).
Permite gerenciar clientes, veículos, produtos e ordens de serviço, com autenticação JWT, documentação via Swagger e deploy simplificado com Docker Compose.

Funcionalidades principais

- Endpoints RESTful organizados por recursos (/customers, /vehicles, /orders, etc.)

- Verificação de saúde (Health Check) em /health

- UI interativa do Swagger em /docs

- Especificação OpenAPI estática servida em /swagger/swagger.yaml

- Autenticação baseada em JWT e middleware de sessão

---

## Visão Geral Rápida

O projeto é uma aplicação Go (module github.com/13SOAT-andromeda/tech-challenge-s1) com ponto de entrada em
cmd/api/main.go.

Por padrão, o servidor escuta na porta 8080 (configurável via variável HTTP_PORT).
A UI do Swagger é servida diretamente pela aplicação e aponta para o arquivo de especificação localizado na pasta swagger/.

---

## Executar com Docker (Recomendado)

O repositório já inclui **Dockerfile** e **docker-compose.yml** para iniciar a aplicação e o banco de dados automaticamente.


Passos

1. Instale o Docker:
    - Windows: https://docs.docker.com/desktop/
    - Ubuntu (exemplo):
      ```bash
      sudo apt update && sudo apt install -y docker.io docker-compose
      sudo systemctl enable --now docker
      ```
2. Crie/ajuste o arquivo .env (mesmo exemplo anterior).
3. Suba os serviços:
   ```bash
   docker-compose up --build
   ```
   ou em background:
   ```bash
   docker-compose up --build -d
   ```
4. Acesse:
    - API: http://localhost:8080/
    - Health: http://localhost:8080/health
    - Swagger UI: http://localhost:8080/docs/
    - Redoc UI: http://localhost:8080/redoc/

Para encerrar:

```bash
docker-compose down
```

---

## Executar localmente (sem Docker)

Pré-requisitos:

- Go 1.25+ instalado
- Uma instância de PostgreSQL em execução (local ou remota) acessível com as credenciais fornecidas

Passos simples:

1. Instale o Go:
    - Windows: baixe e execute o instalador em https://go.dev/dl/
    - Ubuntu (WSL):
      ```bash
        sudo apt update && sudo apt install -y golang
      ```
   2. Crie um arquivo .env na raiz do projeto:

      ```env

       DB_HOST=db
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

      Notas:
      - As variáveis são carregadas em `internal/adapter/config/config.go`.
      - Caso não tenha o PostgreSQL instalado, você pode rodar um contêiner temporário (veja abaixo).

3. Execute a aplicação:
   - Execute diretamente (use `go run` com reconhecimento de módulo):
     ```bash
     go run ./cmd/api
     ```

   - ou:
     ```bash
     go build -o bin/app ./cmd/api
     ./bin/app
     ```

4. Verifique se a aplicação está em execução:
   - Health: http://localhost:8080/health
   - Swagger UI: http://localhost:8080/docs/index.html
   - Redoc: http://localhost:8080/redoc

Banco de dados rápido com Docker
```bash
docker run --name tech-challenge-pg \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres123 \
  -e POSTGRES_DB=myapp_db \
  -p 5432:5432 -d postgres:15
```

Para parar e remover: 
```bash
`docker stop tech-challenge-pg && docker rm tech-challenge-pg`
```
---

## Testes End-to-End (E2E)

Os testes E2E (End-to-End) validam o comportamento completo da aplicação — desde as requisições HTTP até a persistência no banco de dados.
Eles exigem que a aplicação e o banco de dados estejam em execução.

Pré-requisitos

A aplicação deve estar rodando localmente conforme a etapa:   [Executar localmente (sem Docker)](#executar-localmente-sem-docker)

Na raiz do projeto, executar o comando: 

```bash
go test ./test/e2e/... -v
```


## Documentação Swagger / OpenAPI

A documentação da API está disponível em:

Swagger:
http://localhost:8080/docs

Redoc:
http://localhost:8080/redoc

---

## Solução de problemas

- Se a aplicação falhar ao conectar-se ao banco de dados, verifique se o seu Postgres está em execução e se os valores do `.env` correspondem.
- Verifique os logs impressos em stdout/stderr para erros de inicialização.
