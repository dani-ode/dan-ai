Tahap 4.1 — Kafka + Outbox

Yang perlu dibuat dulu hanya ini.

apps/
├── worker-embedding/
│   ├── main.go
│   └── bootstrap/
│       └── worker.go
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
    Aggregate   string
    AggregateID string
    EventType   string
    Timestamp   time.Time
}
Outbox Publisher

Buat module baru.

internal/
└── outbox/
    ├── entity/
    │     outbox.go
    │
    ├── repository/
    │     postgres.go
    │
    ├── service/
    │     service.go
    │
    └── publisher/
          publisher.go

Ini khusus membaca

outbox_events

lalu publish ke Kafka.

Tahap 4.2 — Knowledge Processor

Baru mulai membuat module Knowledge.

internal/
└── knowledge/
    │
    ├── entity/
    │     knowledge_document.go
    │     knowledge_chunk.go
    │
    ├── repository/
    │     postgres.go
    │
    ├── service/
    │     service.go
    │
    ├── builder/
    │     profile.go
    │     project.go
    │     experience.go
    │     certificate.go
    │
    └── processor/
          processor.go

Belum ada AI.

Belum ada Gemini.

Processor hanya

Load Document

↓

Print Document
Tahap 4.3 — AI

Baru setelah Worker stabil.

Tambah folder

internal/
└── ai/
    │
    ├── provider/
    │     gemini.go
    │
    ├── prompt/
    │     knowledge_builder.go
    │
    ├── schema/
    │     knowledge_builder.go
    │
    └── client/
          client.go

Kenapa dipisah?

Karena nanti

Gemini

OpenAI

Claude

Ollama

semuanya tinggal implement interface yang sama.

Tahap 4.4 — Chunk Builder

Sekarang baru.

internal/
└── knowledge/
    └── chunk/
          ai_builder.go

Misalnya

Build(ctx, document)

return

[]Chunk
Tahap 4.5 — Embedding

Setelah chunk jadi.

internal/
└── knowledge/
    └── embedding/
          service.go

Misalnya

Generate(text string)

↓

[]float32
Tahap 4.6 — Milvus

Terakhir.

pkg/
└── milvus/
    ├── client.go
    ├── collection.go
    └── vector.go
Akhirnya struktur menjadi
portfolio-ai/

apps/
├── api/
├── worker-embedding/
└── worker-events/

internal/
├── profile/
├── project/
├── experience/
├── certificate/
├── prompt/
├── visitor/
├── chat/
├── aimodel/
│
├── outbox/
│   ├── entity/
│   ├── repository/
│   ├── service/
│   └── publisher/
│
├── knowledge/
│   ├── entity/
│   ├── repository/
│   ├── service/
│   ├── builder/
│   ├── processor/
│   ├── chunk/
│   └── embedding/
│
├── ai/
│   ├── client/
│   ├── provider/
│   ├── prompt/
│   └── schema/
│
└── shared/

pkg/
├── kafka/
├── milvus/
├── postgres/
├── grpc/
└── logger/
Urutan implementasi yang saya rekomendasikan

Saya akan mengerjakannya persis dalam urutan berikut agar setiap langkah bisa diuji secara independen:

pkg/kafka — koneksi producer dan consumer.
internal/outbox — publisher yang membaca outbox_events lalu mengirim event ke Kafka.
apps/worker-embedding — consumer Kafka yang hanya menerima event dan mencetak log.
internal/knowledge/processor — worker memuat knowledge_document berdasarkan event.
internal/ai — integrasi Gemini/OpenAI beserta prompt dan schema.
internal/knowledge/chunk — mengubah satu knowledge_document menjadi kumpulan self-contained chunks menggunakan LLM.
internal/knowledge/embedding — menghasilkan embedding untuk setiap chunk.
pkg/milvus — menyimpan dan memperbarui vector di Milvus.