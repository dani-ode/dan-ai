Phase 1
────────────────────
Foundation

• Docker
• Gin
• PostgreSQL
• Migrate
• Config
• Logger
• ULID

↓

Phase 2
────────────────────
Core Application

• Authentication
• profile
• prompts
• visitor
• chat_sessions
• chat_messages

↓

Phase 3
────────────────
Knowledge Base

Knowledge Documents
Knowledge Chunks

↓

Phase 4
────────────────
Background Processing

Kafka
Outbox
Event Worker
Embedding Worker

↓

Phase 5
────────────────
Vector Search

Chunking
Embedding
Milvus

↓

Phase 6
────────────────
AI Agent

LLM Provider
AI Models (aktif digunakan)
Prompt Builder
Conversation Memory
RAG
Tool Calling (ai_tools mulai digunakan)

↓

Phase 7
────────────────
Realtime

gRPC Server Streaming
Typing Animation