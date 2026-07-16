FROM golang:alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/worker-memory ./apps/worker-memory

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /bin/worker-memory /bin/worker-memory
ENTRYPOINT ["/bin/worker-memory"]
