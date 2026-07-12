FROM golang:alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/worker-embedding ./apps/worker-embedding

FROM alpine:3.20
COPY --from=build /bin/worker-embedding /bin/worker-embedding
ENTRYPOINT ["/bin/worker-embedding"]
