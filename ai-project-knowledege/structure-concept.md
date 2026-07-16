                        Internet
                            │
          ┌─────────────────┴─────────────────┐
          │                                   │
      gRPCUI                             Vue Frontend
          │                                   │
          └─────────────────┬─────────────────┘
                            │
                       gRPC Server
                            │
                    Chat Service
                            │
          ┌─────────────────┴──────────────────┐
          │                                    │
          ▼                                    ▼
     PostgreSQL                        Kafka Producer
          │                                    │
          │                             chat.completed
          │                                    │
          └────────────────────────────────────┘
                           │
                  Kafka Cluster
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼

Knowledge Worker Memory Worker Event Worker
│ │ │
▼ ▼ ▼
Milvus Visitor Memory Analytics
