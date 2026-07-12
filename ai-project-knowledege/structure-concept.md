                    Internet
                        │
        ┌───────────────┴───────────────┐
        │                               │
   gRPCUI (Admin)                  Vue Frontend
        │                               │
        └───────────────┬───────────────┘
                        │
                  gRPC Server
              (Modular Monolith)
                        │
        ┌───────────────┼───────────────────────────┐
        │               │                           │
        ▼               ▼                           ▼
   PostgreSQL      Kafka Producer             AI Module
(Source of Truth)       │                 (RAG, Prompt Builder,
                        │                  Tool Calling)
                        ▼                           │
                    Kafka Cluster                   │
                  ┌───────────────┐                 │
                  │               │                 │
                  ▼               ▼                 │
          Event Worker   Embedding Worker           │
                                  │                 │
                                  ▼                 │
                               Milvus ◄─────────────┘
                                  │
                                  ▼
                       Semantic Search (Top-K)
                                  │
                                  ▼
                     Gemini / OpenAI / Ollama