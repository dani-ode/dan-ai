1. gRPC Server (Modular Monolith)

Ini adalah satu aplikasi Go yang menjalankan satu gRPC Server.

portfolio-ai

Di dalamnya terdapat module:

Authentication
Profile
Experience
Project
Certificate
Technology
Skill
Prompt
Visitor
Chat
Knowledge
AI
Shared

Semua module berjalan dalam satu process dan dikompilasi menjadi satu binary.

apps/api

berisi:

gRPC Server

Handler

Service

Repository

Middleware

Interceptor

Validator

Mapper

DTO

Semuanya masih berada dalam satu aplikasi Go.

2. PostgreSQL

PostgreSQL adalah source of truth.

Semua module membaca dan menulis data ke PostgreSQL.

Profile
        ↓
profiles

Project
        ↓
projects

Chat
        ↓
chat_messages

Knowledge
        ↓
knowledge_documents

Tidak ada data utama yang disimpan di Kafka maupun Milvus.

Milvus hanyalah index vector.

Kafka hanyalah message broker.

3. Authentication

Karena tidak ada Admin Panel, maka hanya ada dua jenis client.

gRPCUI (Admin)

Vue (Visitor)

Admin menggunakan

Login()

↓

JWT

↓

Authorization:
Bearer xxxx

JWT dipakai untuk semua CRUD.

Misalnya

CreateProject

UpdateProfile

DeleteCertificate

CreatePrompt

Sedangkan visitor tidak memerlukan login.

Frontend cukup mengirim

visitor_id

yang disimpan di LocalStorage.

4. Kafka

Kafka bukan database.

Kafka hanya membawa event.

Contoh.

Admin mengubah Project melalui gRPC.

gRPC

↓

UpdateProject()

↓

UPDATE PostgreSQL

↓

Publish Event

↓

Return Success

Embedding dikerjakan belakangan.

API tidak pernah menunggu embedding selesai.

Topic

Saya akan menyederhanakan topic menjadi berdasarkan domain.

portfolio.profile

portfolio.project

portfolio.certificate

portfolio.knowledge

portfolio.chat

portfolio.embedding

Event Type berada di payload.

Misalnya

{
  "aggregate": "project",
  "aggregate_id": "01K...",
  "event": "updated",
  "timestamp": "...",
  "payload": {}
}

Daripada membuat banyak topic seperti

project.created

project.updated

project.deleted

lebih mudah satu topic per aggregate.

5. Embedding Worker

Worker ini hanya mempunyai satu tanggung jawab.

Mengubah knowledge menjadi embedding.

Misalnya.

Update Project

↓

PostgreSQL

↓

Kafka

↓

Embedding Worker

Worker kemudian.

Ambil Project

↓

Build Document

↓

Chunking

↓

Embedding

↓

Milvus

Worker tidak mempunyai HTTP maupun gRPC Server.

Dia hanya Kafka Consumer.

Semua yang masuk embedding.

Profile

Experience

Project

Certificate

Manual Knowledge

akan diubah menjadi

knowledge_documents

↓

knowledge_chunks

↓

Milvus
6. Event Worker

Worker ini menangani background task selain embedding.

Contohnya.

Update Project

↓

Ambil GitHub API

↓

Update github_stars

atau

Visitor membuka website

↓

Analytics

atau

Certificate dibuat

↓

Generate thumbnail

Embedding Worker dan Event Worker dipisahkan supaya tanggung jawabnya tetap sederhana.

7. Milvus

Milvus bukan database utama.

Milvus hanya menyimpan vector.

Misalnya.

Saya membuat backend menggunakan
Gin,
Kafka,
Docker.

↓

Embedding

↓

[0.21,
0.54,
0.88,
...]

↓

Milvus.

Milvus menyimpan

chunk_id

embedding

metadata

Contoh metadata.

{
    "document_id":"...",
    "source":"project",
    "source_id":"01K..."
}

Detail project tetap berada di PostgreSQL.

8. AI Module

AI bukan service terpisah.

AI cukup menjadi module.

internal/

    ai/

        agent/

        rag/

        prompt/

        providers/

        tools/

        memory/

Semuanya masih berada dalam aplikasi gRPC.

Alur ketika Visitor Chat

Frontend mengirim

SendMessage

session_id

content

Backend melakukan.

Cari Session

↓

Cari Visitor

↓

Cari Prompt

↓

Simpan User Message

↓

Ambil History
(session_id)

↓

Semantic Search
(Milvus)

↓

Build Prompt

↓

LLM

↓

Simpan Assistant Message

↓

Streaming Response

History selalu berdasarkan

session_id

bukan

visitor_id

karena satu visitor bisa mempunyai banyak percakapan.

Alur ketika Admin mengubah Project

Admin menggunakan gRPCUI.

Login

↓

JWT

↓

UpdateProject()

↓

UPDATE projects

↓

UPDATE knowledge_documents

↓

Publish Kafka Event

↓

Return Success

Worker kemudian.

Kafka

↓

Embedding Worker

↓

Chunking

↓

Embedding

↓

Milvus

Admin tidak perlu menunggu proses embedding selesai.