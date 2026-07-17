# Arsitektur Embedding yang Fleksibel

Jangan mendesain Milvus berdasarkan model embedding yang digunakan, tetapi berdasarkan Collection Alias dan Embedding Profile. Dengan cara ini, model embedding dapat diganti kapan saja tanpa mengubah business logic aplikasi.

# Layer 1 — PostgreSQL (Source of Truth)

Semua data asli disimpan di PostgreSQL.

Knowledge
knowledge_documents

knowledge_chunks
----------------
id
document_id
content
created_at
updated_at
Visitor
visitors

chat_sessions

chat_messages

visitor_knowledge
-----------------
id
visitor_id
category
memory_text
importance
created_at
updated_at

PostgreSQL menjadi Source of Truth.

Milvus hanya menyimpan vector untuk kebutuhan semantic search.

# Layer 2 — Embedding Profile

Tambahkan tabel untuk menyimpan konfigurasi embedding yang aktif.

embedding_profiles

id
name
provider
model
dimension
metric_type
knowledge_collection
visitor_collection
is_active
created_at
updated_at

Contoh:

name	provider	model	dimension	knowledge_collection	visitor_collection
e5	Ollama	multilingual-e5-base	768	dan_knowledge_e5	visitor_knowledge_e5
openai-small	OpenAI	text-embedding-3-small	1536	dan_knowledge_openai	visitor_knowledge_openai
bge	Ollama	bge-m3	1024	dan_knowledge_bge	visitor_knowledge_bge

Hanya ada satu profile yang aktif.

is_active = true


# Layer 3 — Milvus Collections

Collection dibuat berdasarkan embedding profile, bukan nama tetap.

Contoh:

dan_knowledge_e5
visitor_knowledge_e5

Jika suatu saat menggunakan OpenAI:

dan_knowledge_openai
visitor_knowledge_openai

Dengan demikian setiap collection memiliki dimensi vector yang sesuai dengan model embedding yang digunakan.

# Layer 4 — Collection Alias

Aplikasi tidak pernah mengetahui nama collection sebenarnya.

Gunakan Alias Milvus.

Contoh:

Alias

dan_knowledge
        │
        ▼
dan_knowledge_e5

Ketika migrasi ke model baru:

Alias

dan_knowledge
        │
        ▼
dan_knowledge_openai

Kode aplikasi tetap:

Search("dan_knowledge")

Tanpa perubahan apa pun.

Hal yang sama berlaku untuk:

visitor_knowledge
        │
        ▼
visitor_knowledge_e5


# Layer 5 — Embedding Service

Jangan mengikat service dengan model tertentu.

Bukan:

GenerateEmbedding(text)

Tetapi:

GenerateEmbedding(ctx, profile, text)

Contoh:

profile := embeddingProfile.Active()

vector, err := embedding.Generate(
    ctx,
    profile,
    text,
)

Service hanya mengetahui Embedding Profile, bukan Gemini, OpenAI, ataupun Ollama.

# Layer 6 — Search Service

Search juga tidak boleh mengetahui nama collection.

Bukan:

SearchKnowledge(vector)

Tetapi:

SearchKnowledge(
    profile,
    vector,
)

Service membaca:

knowledge_collection
visitor_collection

dari embedding_profiles, kemudian melakukan pencarian menggunakan alias yang sesuai.


# Mengganti Model Embedding

Misalnya saat ini menggunakan:

Model

multilingual-e5-base

Dimension

768

Collection

dan_knowledge_e5
visitor_knowledge_e5

Lalu ingin berpindah ke OpenAI.

Langkahnya:

Buat Embedding Profile baru
        │
        ▼
Buat Collection baru
        │
        ▼
Worker melakukan Re-index
        │
        ▼
Update Alias Milvus
        │
        ▼
Selesai

Tidak perlu mengubah kode Chat Service.

Tidak perlu restart API.

# Embedding Worker

Worker selalu membaca profile yang aktif.

Embedding Profile
        │
        ▼
Generate Embedding
        │
        ▼
Insert / Update Vector
        │
        ▼
Milvus Collection

Worker tidak pernah melakukan hardcode seperti:

"dan_knowledge"

# Embedding Provider

Setiap provider mengimplementasikan interface yang sama.

type EmbeddingProvider interface {
    Name() string
    Dimension() int

    Generate(
        ctx context.Context,
        text string,
    ) ([]float32, error)
}

Implementasi:

GeminiProvider

OpenAIProvider

OllamaProvider

BGEProvider

Sehingga pergantian provider tidak memengaruhi business logic aplikasi.


# Arsitektur Akhir
PostgreSQL (Source of Truth)
────────────────────────────────

knowledge_documents

knowledge_chunks

visitor_knowledge

embedding_profiles

            │
            ▼

Embedding Profile
────────────────────────────────

Model
Dimension
Collections

            │
            ▼

Embedding Service
────────────────────────────────

Gemini
OpenAI
Ollama
BGE

            │
            ▼

Milvus (Vector Index)
────────────────────────────────

Alias
dan_knowledge
        │
        ▼
dan_knowledge_e5

Alias
visitor_knowledge
        │
        ▼
visitor_knowledge_e5

# Struktur visitor_knowledge
## PostgreSQL
visitor_knowledge

id
visitor_id
category
memory_text
importance
created_at
updated_at

## Milvus
visitor_knowledge_<profile>

id              // sama dengan PostgreSQL.id
visitor_id
embedding

---- Dengan desain ini:

PostgreSQL menyimpan seluruh isi memory sebagai Source of Truth.
Milvus hanya menyimpan vector index.
Pergantian model embedding cukup dengan membuat collection baru, melakukan re-index, lalu mengubah alias Milvus tanpa mengubah kode aplikasi maupun database utama.