<!-- re run the app -->
docker compose --env-file .env -f deployments/compose/docker-compose.yml down

$env:GOOS='linux'; $env:GOARCH='amd64'; $env:CGO_ENABLED='0'; go build -o bin/api ./apps/api

docker compose --env-file .env -f deployments/compose/docker-compose.yml up --build -d


<!-- re run the app and the refresh the migration -->
docker compose --env-file .env -f deployments/compose/docker-compose.yml down -v
docker compose --env-file .env -f deployments/compose/docker-compose.yml --profile tools run --rm migrate -path=/migrations -database="postgres://postgres:postgres@postgres:5432/portfolio_ai?sslmode=disable" up