# dan-ai

Dan AI backend built with Go and Gin. The project includes a main API service, background workers for embeddings and events, and supporting infrastructure such as PostgreSQL, Kafka, and Milvus.

## Project layout

- `apps/api`: main HTTP API
- `apps/worker-embedding`: embedding worker
- `apps/worker-events`: event worker
- `internal`: domain and shared application logic
- `pkg`: reusable infrastructure helpers
- `deployments`: Docker and Compose configuration

## Prerequisites

- Docker and Docker Compose
- Go 1.25+
- Make (optional)

## Environment setup

1. Copy `.env.example` to `.env`.
2. Review the values in `.env` before starting the services.

## Run locally

### Start the full stack

```bash
docker compose --env-file .env -f deployments/compose/docker-compose.yml up --build -d
```

### Rebuild and restart the API

```bash
docker compose --env-file .env -f deployments/compose/docker-compose.yml down
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/api ./apps/api
docker compose --env-file .env -f deployments/compose/docker-compose.yml up --build -d
```

### Reset the database and rerun migrations

```bash
docker compose --env-file .env -f deployments/compose/docker-compose.yml down -v
docker compose --env-file .env -f deployments/compose/docker-compose.yml --profile tools run --rm migrate -path=/migrations -database="postgres://postgres:postgres@postgres:5432/dan_ai?sslmode=disable" up
```

## Verify services

```bash
docker compose --env-file .env -f deployments/compose/docker-compose.yml ps
```

You should see the main services running normally, including:

- `dan-ai-api` on port `8080`
- `dan-ai-postgres` on port `5432`
- `dan-ai-kafka` on ports `9092/9093`
- `dan-ai-milvus` on port `19530`
