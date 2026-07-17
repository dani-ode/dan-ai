<!-- re run the app -->

make docker-up

<!-- or -------------------------------------------- -->

docker compose --env-file .env -f deployments/compose/docker-compose.yml down

$env:GOOS='linux'; $env:GOARCH='amd64'; $env:CGO_ENABLED='0'; go build -o bin/api ./apps/api

# Optional: build worker binaries locally for testing

# On Linux/macOS:

# export GOOS=linux; export GOARCH=amd64; export CGO_ENABLED=0; go build -o bin/worker-memory ./apps/worker-memory

docker compose --env-file .env -f deployments/compose/docker-compose.yml up --build -d

<!-- re run the app and the refresh the migration -->

docker compose --env-file .env -f deployments/compose/docker-compose.yml down -v
docker compose --env-file .env -f deployments/compose/docker-compose.yml --profile tools run --rm migrate -path=/migrations -database="postgres://postgres:postgres@postgres:5432/dan_ai?sslmode=disable" up

<!-- the result must be -->

All containers are running normally:

dan-ai-api (Up, on port 8080)
dan-ai-postgres (Up, Healthy, on port 5432)
dan-ai-kafka (Up, on ports 9092/9093)
dan-ai-milvus (Up, on port 19530)
dan-ai-milvus-minio, dan-ai-milvus-etcd, dan-ai-worker-embedding, dan-ai-worker-events, and dan-ai-worker-memory are all running and active.
