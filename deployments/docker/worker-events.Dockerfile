FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /bin/worker-events ./apps/worker-events

FROM alpine:3.20
COPY --from=build /bin/worker-events /bin/worker-events
ENTRYPOINT ["/bin/worker-events"]
