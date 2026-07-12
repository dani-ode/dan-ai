# deployments/docker/api.Dockerfile

FROM golang:alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./apps/api

FROM alpine:3.20
COPY --from=build /bin/api /bin/api
ENTRYPOINT ["/bin/api"]