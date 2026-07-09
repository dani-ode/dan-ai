# deployments/docker/api.Dockerfile

FROM alpine:3.21

WORKDIR /app

COPY bin/api /bin/api

EXPOSE 8080

ENTRYPOINT ["/bin/api"]