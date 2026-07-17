# Container Logging Guide

This guide explains how to view, manage, and monitor the logs of the containerized services in the **Dan AI** platform.

Since the entire application runs inside Docker Compose, all standard output (`stdout`) and standard error (`stderr`) logs are captured automatically. You do **not** need to create or write to separate `.log` files manually.

---

## 1. Monitoring All Services (Real-time)

To view a combined live stream of logs from all active containers:

```bash
docker compose --env-file .env -f deployments/compose/docker-compose.yml logs -f
```

---

## 2. Inspecting Individual Services

If you are debugging a specific component of the platform, filter the logs by the service name:

### AI Embedding Worker
Responsible for consuming events from Kafka, chunking documents using Gemini, generating embeddings, and indexing them to Milvus.
```bash
docker compose --env-file .env -f deployments/compose/docker-compose.yml logs worker-embedding
```

### Memory Consolidation Worker
Responsible for parsing chat histories, calling Gemini to extract visitor context, searching Milvus, and merging highly similar memories.
```bash
docker compose --env-file .env -f deployments/compose/docker-compose.yml logs worker-memory
```

### API Server
Handles all incoming gRPC and HTTP client requests, GORM database connections, and session registrations.
```bash
docker compose --env-file .env -f deployments/compose/docker-compose.yml logs api
```

### Outbox Event Worker
Polls the PostgreSQL `outbox_events` table and publishes unpublished events to the respective Kafka topics (`dan.knowledge` or `dan.events`).
```bash
docker compose --env-file .env -f deployments/compose/docker-compose.yml logs worker-events
```

---

## 3. General Troubleshooting Tips

### Checking Container Statuses
To check if any container is restarting or unhealthy:
```bash
docker compose --env-file .env -f deployments/compose/docker-compose.yml ps
```

### Clearing/Restarting Logs
If the log output becomes too long and you want to start fresh:
```bash
docker compose --env-file .env -f deployments/compose/docker-compose.yml restart <service-name>
```
