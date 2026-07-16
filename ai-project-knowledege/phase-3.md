# Langkah 1 — AI Models

Ini paling mudah.

Karena nanti AI Agent harus mengetahui model apa yang aktif.

Misalnya

Gemini 2.5 Flash

temperature = 0.2

context = 32000

enabled = true

CRUD biasa saja.

Nanti Phase 6 tinggal

SELECT * FROM ai_models
WHERE enabled = true
LIMIT 1


# Langkah 2 — Knowledge Document

Ini mulai menarik.

Misalnya ada Project.

Project

ID
01ABC

Title
Dan AI

Summary
AI Portfolio

Description
Backend menggunakan Golang...

Belum bisa langsung dikirim ke embedding.

Maka dibuat satu document.

knowledge_documents

id
01XYZ

source_type
project

source_id
01ABC

title
Dan AI

content
Dan AI adalah backend...

Yang mengisi content bukan AI.

Melainkan Builder.

Misalnya

func BuildProjectDocument(project Project) string

menghasilkan

Project:
Dan AI

Summary:
Backend AI Portfolio.

Architecture:
Kafka
Milvus

Description:
Dan AI merupakan backend...

Kalimatnya dibuat deterministic.

Tidak perlu AI.

Jadi nanti setiap entity punya Builder
knowledge/
    builder/

        project.go

        experience.go

        certificate.go

        profile.go

        skill.go

Misalnya

BuildProjectDocument()

BuildProfileDocument()

BuildCertificateDocument()


# Langkah 3 — Generate Knowledge Document

Sekarang CRUD Project sudah ada.

Sesudah

CreateProject()

langsung

BuildProjectDocument()

↓

INSERT knowledge_documents

Begitu juga

UpdateProject()

↓

Update knowledge_document

Jadi belum ada Kafka dulu.

Masih synchronous.

Nanti saat Kafka masuk tinggal dipindah ke Worker.

# Langkah 4 — Knowledge Chunks

Sesudah ada document.

Misalnya

Project

2000 kata

AI embedding tidak boleh sepanjang itu.

Harus dipotong.

Misalnya

Chunk 1

Dan AI adalah...

--------------------

Chunk 2

Project ini menggunakan...

--------------------

Chunk 3

Deployment dilakukan...

Masuk ke

knowledge_chunks

Isinya

id

document_id

chunk_index

content

token_count

embedding_model

Belum ada embedding.

Masih text.

# Langkah 5 — Chunk Builder

Misalnya

SplitEvery(500 tokens)

atau

SplitEvery(800 characters)

sementara.

Nanti Phase 5 baru pakai tokenizer.

Folder

knowledge/

    chunker/

        chunker.go

Contoh

BuildChunks(document)

hasil

[]

Chunk

Chunk

Chunk

langsung

INSERT knowledge_chunks
Langkah 6 — Repository Knowledge

Buat CRUD sederhana.

KnowledgeDocumentRepository

Create()

Update()

Delete()

GetBySource()

List()

dan

KnowledgeChunkRepository

CreateMany()

DeleteByDocument()

ListByDocument()
Langkah 7 — Service Knowledge

Service bertugas mengorkestrasi semuanya.

Misalnya

Project dibuat

↓

ProjectService.Create()

↓

KnowledgeService.SyncProject()

↓

Build Document

↓

Save Document

↓

Chunking

↓

Save Chunks

Diagramnya

Project

↓

Knowledge Builder

↓

Knowledge Document

↓

Chunk Builder

↓

Knowledge Chunks
Struktur folder
internal/

    knowledge/

        entity/

            document.go

            chunk.go

        builder/

            profile.go

            project.go

            experience.go

            certificate.go

            skill.go

        chunker/

            chunker.go

        repository/

            postgres.go

        service/

            service.go

        grpc/

            handler.go

        mapper/

            mapper.go


# Hasil akhir Phase 3

Setelah Phase 3 selesai, alurnya menjadi seperti ini:

Create Project

        │

        ▼

Save Project
(PostgreSQL)

        │

        ▼

Generate Knowledge Document

        │

        ▼

Save knowledge_documents

        │

        ▼

Split menjadi Chunks

        │

        ▼

Save knowledge_chunks

Belum ada:

❌ Kafka
❌ Embedding
❌ Milvus
❌ Gemini/OpenAI
❌ RAG

Namun ketika masuk Phase 4, worker embedding tinggal mengambil data dari knowledge_chunks untuk membuat embedding tanpa perlu mengubah modul Project, Experience, atau Certificate lagi. Inilah keuntungan memisahkan Knowledge Layer sebagai fondasi sebelum masuk ke pipeline AI.