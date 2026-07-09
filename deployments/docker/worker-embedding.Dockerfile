FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /bin/worker-embedding ./apps/worker-embedding

FROM alpine:3.20
COPY --from=build /bin/worker-embedding /bin/worker-embedding
ENTRYPOINT ["/bin/worker-embedding"]
