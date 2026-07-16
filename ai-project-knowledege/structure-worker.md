# Worker dibagi menjadi 3

## 1. Knowledge Worker

Hanya menangani:

Profile Updated
Project Updated
Experience Updated
Certificate Updated

↓

Build Document

↓

Chunk

↓

Embedding

↓

Milvus (knowledge)

Collection:

knowledge

## 2. Memory Worker

Hanya menangani

chat.completed

↓

Load Chat

↓

Extract Memory

↓

Embedding

↓

Milvus (visitor_memory)

Collection:

visitor_memory

## 3. Event Worker

Tidak berhubungan dengan AI.

Misalnya

audit log

analytics

counter

visitor statistic

daily report

# Struktur folder

Aku akan mengubah sedikit.

apps/

api/

worker-knowledge/

worker-memory/

worker-events/

Lalu di internal

internal/

knowledge/

memory/

chat/

visitor/

Tambahkan

internal/memory/

entity/

repository/

service/

extractor/

prompts/

# Chat Flow

Misalnya user bertanya.

User

↓

ChatService

ChatService hanya melakukan

Save User Message

↓

Search Knowledge

↓

Search Visitor Memory

↓

Build Prompt

↓

LLM

↓

Save Assistant

↓

Publish chat.completed

Selesai.

Tidak ada embedding.

Tidak ada memory extraction.

Supaya response cepat.

# Memory Worker

Worker menerima

chat.completed

Payload

{
"visitor_id":"01K...",
"session_id":"01K...",
"user_message_id":"01K...",
"assistant_message_id":"01K..."
}

Perhatikan:

Aku tidak mengirim isi chat.

Worker tinggal query PostgreSQL.

SELECT
user_message,
assistant_message

Lalu prompt

Extract useful memory.

Ignore greetings.

Ignore thanks.

Ignore generic questions.

Return JSON.

Misalnya

{
"save": true,
"category":"project",
"importance":5,
"memory":"Visitor sedang membangun Portfolio AI menggunakan Kafka."
}

Kalau

{
"save":false
}

langsung selesai.

Kalau save=true

Memory

↓

Embedding

↓

Milvus Search
(filter visitor_id)

↓

Similarity > 0.9 ?

Kalau ada

Merge

↓

Update

Kalau tidak ada

Insert

# Collection Milvus

Aku justru akan membuat sederhana.

## dan_knowledge

id

document_id

source_type

source_id

text

embedding

metadata

## visitor_memory

id

visitor_id

session_id

category

importance

memory

embedding

created_at

Tidak perlu chunk.

Memory maksimal 250 karakter.

# PostgreSQL

Tetap menjadi source of truth.

visitors

chat_sessions

chat_messages

visitor_memories

Kalau visitor_memories ingin dipakai, isinya adalah fakta terstruktur seperti:

favorite_language = Golang

city = Surabaya

current_project = Portfolio AI

Sedangkan Milvus menyimpan memori semantik.
