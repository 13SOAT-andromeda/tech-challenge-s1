# Tech Challenge S1

A small Go HTTP API project that exposes a few customer-related endpoints and serves OpenAPI (Swagger) documentation.

Key features
- Simple REST endpoints under `/customer`.
- Health check at `/health`.
- Serves Swagger UI at `/docs` and the static OpenAPI spec at `/swagger/swagger.yaml`.

---

## Quick overview
This project is a Go application (module `github.com/13SOAT-andromeda/tech-challenge-s1`) with the server entrypoint in `cmd/api/main.go`. By default it listens on port `8080` (configurable via `HTTP_PORT`).

The Swagger UI is served by the application and points to the static spec file shipped in the `swagger/` folder.

---

## Run locally (without Docker)
Prerequisites:
- Go 1.25+ installed
- A running PostgreSQL instance (local or remote) reachable with the credentials you provide

Simple steps:
1. Install Go (one-liner examples):
   - Windows: download and run the installer from https://go.dev/dl/
   - Ubuntu (WSL):
     ```bash
     sudo apt update && sudo apt install -y golang
     ```
2. From the project root, create a `.env` file or set environment variables. Example `.env` (place in project root):

   ```env
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=password
   DB_NAME=postgres
   DB_PORT=5432
   DB_SSLMODE=disable
   DB_TIMEZONE=UTC

   HTTP_PORT=8080
   HTTP_URL=http://localhost
   HTTP_ALLOWED_ORIGINS=*
   ENV=development
   ```

   Notes:
   - The app reads environment variables defined in `internal/adapter/config/config.go`.
   - If you don't want to install Postgres locally, you can temporarily start one using Docker (example below), but that uses Docker.

3. Run the app:
   - Run directly (use module-aware `go run`):
     ```bash
     go run ./cmd/api
     ```

   - Or build and run a binary:
     ```bash
     go build -o bin/app ./cmd/api
     ./bin/app
     ```

4. Verify the app is running:
   - Health: http://localhost:8080/health
   - Swagger UI: http://localhost:8080/docs/index.html?url=/swagger/swagger.yaml

Optional: quick Postgres (Docker) for local testing:
```bash
docker run --name tech-challenge-pg -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=postgres -p 5432:5432 -d postgres:15
```
(Stop/remove it when done with `docker stop tech-challenge-pg && docker rm tech-challenge-pg`)

---

## Run with Docker (recommended for a quick isolated setup)
This repository includes a `Dockerfile` and a `docker-compose.yml` that start a Postgres DB and the application.

1. Install Docker Desktop (Windows) or Docker Engine + Compose (Linux). Simple pointers:
   - Docker Desktop for Windows: https://docs.docker.com/desktop/
   - Ubuntu (example):
     ```bash
     sudo apt update && sudo apt install -y docker.io docker-compose
     sudo systemctl enable --now docker
     ```
2. From the project root, create or edit `.env` with DB credentials (same as above).
3. Start services with compose:
   ```bash
   docker-compose up --build
   ```
   Or run in background:
   ```bash
   docker-compose up --build -d
   ```
4. Endpoints:
   - App: http://localhost:8080/
   - Health: http://localhost:8080/health
   - Swagger UI: http://localhost:8080/docs/index.html?url=/swagger/swagger.yaml

To stop and remove containers:
```bash
docker-compose down
```

---

## Swagger / OpenAPI docs
The Swagger UI is available once the app is running. Open this URL in your browser:

http://localhost:8080/docs/index.html?url=/swagger/swagger.yaml

This points the Swagger UI to the static spec file served at `/swagger/swagger.yaml`.

---

## Example API endpoints
- GET /health -> health check
- GET /customer -> list customers
- POST /customer -> create a customer
- GET /customer/:id -> get customer by id

---

## Troubleshooting
- If the app fails to connect to the database, verify your Postgres is running and the `.env` values match.
- Check logs printed to stdout/stderr for startup errors.

---

If you'd like, I can add a minimal Makefile or scripts to simplify run commands, or add example curl requests for the endpoints.
