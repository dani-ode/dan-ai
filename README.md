# portfolio-ai

Go + Gin scaffold for the portfolio-ai backend.

## Layout

- `apps/api`: main HTTP API
- `apps/worker-embedding`: embedding worker
- `apps/worker-events`: event worker
- `internal`: domain and shared application code
- `pkg`: reusable infrastructure helpers
- `deployments`: Docker and compose files

## Run

1. Copy `.env.example` to `.env`.
2. Run the API with `go run ./apps/api`.
