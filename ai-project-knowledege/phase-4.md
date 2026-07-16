# Tahap 4.1 — Kafka + Outbox

Yang dibuat terlebih dahulu adalah infrastruktur event.

apps/
├── worker-knowledge/
│ ├── main.go
│ └── bootstrap/
│ └── worker.go
│
├── worker-memory/
│ ├── main.go
│ └── bootstrap/
│ └── worker.go
│
└── worker-events/
├── main.go
└── bootstrap/
└── worker.go

Kemudian di pkg/

pkg/
└── kafka/
├── producer.go
├── consumer.go
├── event.go
└── kafka.go

event.go

type Event struct {
Aggregate string
AggregateID string
EventType string
Timestamp time.Time
}

# Outbox Publisher

Buat module khusus untuk membaca outbox_events dan mengirimkannya ke Kafka.

internal/
└── outbox/
├── entity/
│ outbox.go
│
├── repository/
│ postgres.go
│
├── service/
│ service.go
│
└── publisher/
publisher.go

Pada tahap ini worker hanya memastikan event berhasil dipublish ke Kafka.

# Tahap 4.2 — Knowledge Processor

Mulai membuat pipeline Knowledge.

internal/
└── knowledge/
│
├── entity/
│ knowledge_document.go
│ knowledge_chunk.go
│
├── repository/
│ postgres.go
│
├── service/
│ service.go
│
├── builder/
│ profile.go
│ project.go
│ experience.go
│ certificate.go
│ skill.go
│
└── processor/
processor.go

Flow sementara:

Kafka Event

↓

Load Document

↓

Build Knowledge Document

↓

Print Result

Belum ada AI.

Belum ada Chunk.

Belum ada Embedding.

# Tahap 4.3 — Memory Processor

Mulai membuat pipeline untuk Long-Term Memory Visitor.

internal/
└── memory/
│
├── entity/
│ memory.go
│
├── repository/
│ postgres.go
│ milvus.go
│
├── service/
│ service.go
│
├── extractor/
│ extractor.go
│
└── processor/
processor.go

Flow sementara:

chat.completed

↓

Load User Message

↓

Load Assistant Message

↓

Print Conversation

Belum ada AI.

Belum ada Embedding.

# Tahap 4.4 — AI Layer

Setelah kedua worker stabil, baru menambahkan AI.

internal/
└── ai/
│
├── client/
│ client.go
│
├── provider/
│ gemini.go
│ openai.go
│ ollama.go
│
├── prompt/
│ knowledge_builder.go
│ memory_extractor.go
│ memory_merge.go
│
└── schema/
knowledge_builder.go
memory_extractor.go

Semua provider mengimplementasikan interface yang sama sehingga mudah diganti.

# Tahap 4.5 — Knowledge Chunk Builder

internal/
└── knowledge/
└── chunk/
ai_builder.go

Flow

Knowledge Document

↓

LLM

↓

Self-contained Chunks

# Tahap 4.6 — Memory Extractor

Sekarang AI mulai membuat Long-Term Memory.

Flow

User Message

-

Assistant Message

↓

LLM

↓

{
save,
category,
importance,
memory
}

Jika save=false, worker selesai.

Jika save=true, lanjut ke embedding.

# Tahap 4.7 — Embedding

Knowledge dan Memory menggunakan service embedding yang sama.

internal/
├── knowledge/
│ └── embedding/
│ service.go
│
└── memory/
└── embedding/
service.go

Contoh:

Generate(text)

↓

[]float32

# Tahap 4.8 — Milvus

Terakhir membuat wrapper Milvus.

pkg/
└── milvus/
├── client.go
├── collection.go
├── vector.go
└── search.go

Collection yang digunakan:

dan_knowledge

Berisi seluruh knowledge tentang dirimu:

Profile
Experience
Project
Certificate
Skill
Technology
Prompt (jika diperlukan)

Contoh schema:

id
document_id
chunk_id
source_type
source_id
title
content
embedding
metadata
created_at
updated_at

Dan collection kedua:

visitor_knowledge

Berisi long-term memory setiap visitor.

Contoh schema:

id
visitor_id
session_id
category
importance
memory
embedding
created_at
updated_at

Contoh isi:

Visitor sedang membangun Portfolio AI menggunakan Kafka.

Visitor lebih menyukai Golang dibanding Laravel.

Visitor tinggal di Surabaya.

Visitor sedang mencari pekerjaan Backend Engineer.
Saat Chat

Flow retrieval menjadi:

User Question
│
▼
Embedding
│
├──────────────────────────┐
▼ ▼
dan_knowledge visitor_knowledge
Top 5 Top 4
(filter visitor_id)
│ │
└──────────────┬───────────┘
▼
Prompt Builder
▼
LLM
Kenapa tidak satu collection?

Secara teknis memang bisa dibuat satu collection:

knowledge_type

dan
visitor

Lalu difilter.

Tetapi saya tidak menyarankan itu karena karakteristik datanya berbeda.

dan_knowledge visitor_knowledge
Data relatif statis Data terus bertambah
Di-update saat CMS berubah Di-update setiap percakapan
Ribuan chunk Bisa puluhan ribu memory
Tidak perlu filter visitor Selalu filter visitor_id

Memisahkan collection membuat konfigurasi index, strategi pembersihan (retention), dan optimasi query bisa berbeda jika suatu saat diperlukan.

# Urutan implementasi yang saya rekomendasikan

Saya akan mengerjakannya dalam urutan berikut agar setiap tahap dapat diuji secara independen:

pkg/kafka — koneksi producer dan consumer.
internal/outbox — publisher yang membaca outbox_events lalu mengirim event ke Kafka.
apps/worker-knowledge — consumer Kafka yang menerima event knowledge dan mencetak log.
apps/worker-memory — consumer Kafka yang menerima event chat.completed dan mencetak log.
internal/knowledge/processor — memuat knowledge_document berdasarkan event.
internal/memory/processor — memuat user_message dan assistant_message berdasarkan event.
internal/ai — integrasi Gemini/OpenAI/Ollama beserta prompt dan schema.
internal/knowledge/chunk — mengubah knowledge_document menjadi kumpulan self-contained chunks menggunakan LLM.
internal/memory/extractor — menghasilkan long-term memory dari satu percakapan.
internal/knowledge/embedding dan internal/memory/embedding — menghasilkan embedding vector.
pkg/milvus — menyimpan, memperbarui, dan melakukan semantic search pada collection knowledge dan visitor_memory.

Dengan urutan ini, pipeline Knowledge RAG dan Conversation Memory berkembang secara paralel, tetapi tetap terpisah tanggung jawabnya. Hal ini membuat kode lebih modular, lebih mudah diuji, dan lebih mudah dikembangkan ketika nanti kamu menambahkan fitur seperti memory consolidation, reranking, atau multi-agent AI.
