FROM golang:alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/worker-events ./apps/worker-events

FROM alpine:3.20
COPY --from=build /bin/worker-events /bin/worker-events
ENTRYPOINT ["/bin/worker-events"]
