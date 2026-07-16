dan-ai/
│
├── .dockerignore
├── .env
├── .env.example
├── .gitignore
├── Makefile
├── README.md
├── apps/
│ ├── api/
│ │ ├── bootstrap/
│ │ │ ├── app.go
│ │ │ ├── database.go
│ │ │ └── grpc.go
│ │ └── main.go
│ ├── worker-embedding/
│ │ ├── bootstrap/
│ │ │ └── worker.go
│ │ └── main.go
│ ├── worker-events/
│ │ ├── bootstrap/
│ │ │ └── worker.go
│ │ └── main.go
│ └── worker-memory/
│ ├── bootstrap/
│ │ └── worker.go
│ └── main.go
├── buf.gen.yaml
├── buf.yaml
├── deployments/
│ ├── compose/
│ │ ├── docker-compose.dev.yml
│ │ ├── docker-compose.prod.yml
│ │ └── docker-compose.yml
│ ├── docker/
│ │ ├── api.Dockerfile
│ │ ├── worker-embedding.Dockerfile
│ │ ├── worker-events.Dockerfile
│ │ └── worker-memory.Dockerfile
│ └── migrations/
│ ├── 000001_create_profiles.down.sql
│ ├── 000001_create_profiles.up.sql
│ ├── 000002_create_remaining_tables.down.sql
│ └── 000002_create_remaining_tables.up.sql
├── docs/
│ └── README.md
├── go.mod
├── go.sum
├── internal/
│ ├── ai/
│ │ ├── client/
│ │ │ └── client.go
│ │ ├── memory/
│ │ ├── provider/
│ │ │ └── gemini.go
│ │ ├── rag/
│ │ ├── repository/
│ │ ├── schema/
│ │ │ └── knowledge_builder.go
│ │ └── service/
│ ├── aimodel/
│ │ ├── entity/
│ │ │ └── aimodel.go
│ │ ├── grpc/
│ │ │ └── handler.go
│ │ ├── mapper/
│ │ │ └── mapper.go
│ │ ├── repository/
│ │ │ └── postgres.go
│ │ └── service/
│ ├── auth/
│ │ ├── grpc/
│ │ │ └── handler.go
│ │ └── jwt/
│ │ └── jwt.go
│ ├── certificate/
│ │ ├── entity/
│ │ │ └── certificate.go
│ │ ├── grpc/
│ │ │ └── handler.go
│ │ ├── mapper/
│ │ │ └── mapper.go
│ │ ├── repository/
│ │ │ └── postgres.go
│ │ └── service/
│ ├── chat/
│ │ ├── entity/
│ │ │ ├── message.go
│ │ │ └── session.go
│ │ ├── grpc/
│ │ │ └── handler.go
│ │ ├── mapper/
│ │ │ └── mapper.go
│ │ ├── repository/
│ │ │ └── postgres.go
│ │ └── service/
│ │ └── service.go
│ ├── experience/
│ │ ├── entity/
│ │ │ └── experience.go
│ │ ├── grpc/
│ │ │ └── handler.go
│ │ ├── mapper/
│ │ │ └── mapper.go
│ │ ├── repository/
│ │ │ └── postgres.go
│ │ └── service/
│ ├── knowledge/
│ │ ├── builder/
│ │ │ ├── certificate.go
│ │ │ ├── experience.go
│ │ │ ├── profile.go
│ │ │ ├── project.go
│ │ │ └── skill.go
│ │ ├── chunk/
│ │ │ └── ai_builder.go
│ │ ├── chunker/
│ │ │ └── chunker.go
│ │ ├── entity/
│ │ │ ├── chunk.go
│ │ │ └── document.go
│ │ ├── grpc/
│ │ │ └── handler.go
│ │ ├── mapper/
│ │ │ └── mapper.go
│ │ ├── processor/
│ │ │ └── processor.go
│ │ ├── repository/
│ │ │ └── postgres.go
│ │ └── service/
│ │ └── service.go
│ ├── memory/
│ │ ├── entity/
│ │ │ └── memory.go
│ │ ├── extractor/
│ │ │ └── extractor.go
│ │ ├── processor/
│ │ │ └── processor.go
│ │ ├── repository/
│ │ │ └── postgres.go
│ │ └── service/
│ │ └── service.go
│ ├── outbox/
│ │ ├── entity/
│ │ │ └── outbox.go
│ │ ├── publisher/
│ │ │ └── publisher.go
│ │ └── repository/
│ │ └── postgres.go
│ ├── profile/
│ ├── project/
│ ├── prompt/
│ ├── shared/
│ │ ├── constants/
│ │ │ └── constants.go
│ │ ├── errors/
│ │ │ └── errors.go
│ │ ├── interceptor/
│ │ │ ├── auth.go
│ │ │ ├── logger.go
│ │ │ └── recovery.go
│ │ └── response/
│ │ └── response.go
│ ├── skill/
│ ├── technology/
│ └── visitor/
├── pkg/
│ ├── config/
│ │ └── config.go
│ ├── grpc/
│ │ └── server.go
│ ├── kafka/
│ │ ├── consumer.go
│ │ ├── event.go
│ │ ├── kafka.go
│ │ └── producer.go
│ ├── logger/
│ │ └── logger.go
│ ├── milvus/
│ │ ├── client.go
│ │ ├── collection.go
│ │ └── vector.go
│ ├── postgres/
│ │ └── postgres.go
│ └── ulid/
│ └── ulid.go
├── proto/
│ ├── aimodel/
│ ├── auth/
│ ├── certificate/
│ ├── chat/
│ ├── experience/
│ ├── knowledge/
│ ├── profile/
│ ├── project/
│ ├── prompt/
│ ├── skill/
│ ├── technology/
│ └── visitor/
├── scripts/
└── README.md
