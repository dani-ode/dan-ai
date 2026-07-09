portfolio-ai/
│
├── apps/
│   ├── api/
│   │   ├── main.go
│   │   └── bootstrap/
│   │       ├── app.go
│   │       ├── grpc.go
│   │       ├── database.go
│   │       └── config.go
│   │
│   ├── worker-embedding/
│   │   ├── main.go
│   │   └── bootstrap/
│   │       └── worker.go
│   │
│   └── worker-events/
│       ├── main.go
│       └── bootstrap/
│           └── worker.go
│
├── internal/
│   │
│   ├── profile/
│   │   ├── entity/
│   │   ├── repository/
│   │   ├── service/
│   │   ├── grpc/
│   │   └── mapper/
│   │
│   ├── visitor/
│   │   ├── entity/
│   │   ├── repository/
│   │   ├── service/
│   │   ├── grpc/
│   │   └── mapper/
│   │
│   ├── chat/
│   │   ├── entity/
│   │   ├── repository/
│   │   ├── service/
│   │   ├── grpc/
│   │   └── mapper/
│   │
│   ├── prompt/
│   │   ├── entity/
│   │   ├── repository/
│   │   ├── service/
│   │   ├── grpc/
│   │   └── mapper/
│   │
│   ├── knowledge/
│   │
│   ├── ai/
│   │
│   └── shared/
│       ├── errors/
│       ├── middleware/
│       ├── validator/
│       └── constants/
│
├── pkg/
│   ├── config/
│   ├── postgres/
│   ├── grpc/
│   ├── logger/
│   ├── ulid/
│   └── utils/
│
├── proto/
│   ├── common/
│   │
│   ├── profile/
│   │
│   ├── visitor/
│   │
│   ├── chat/
│   │
│   ├── prompt/
│   │
│   ├── knowledge/
│   │
│   └── ai/
│
├── deployments/
│   ├── compose/
│   │   ├── docker-compose.yml
│   │   ├── docker-compose.dev.yml
│   │   └── docker-compose.prod.yml
│   │
│   ├── docker/
│   │   ├── api.Dockerfile
│   │   ├── worker-embedding.Dockerfile
│   │   └── worker-events.Dockerfile
│   │
│   └── migrations/
│
├── docs/
│
├── scripts/
│
├── .env
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md