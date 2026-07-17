Visitor
   │
   ▼
Frontend (Vue)
   │
   ▼
ChatService
   │
   ├── Save User Message (PostgreSQL)
   │
   ├── Generate Query Embedding
   │
   ├── Search dan_knowledge (Milvus)
   │
   ├── Search visitor_knowledge (Milvus)
   │      (filter visitor_id)
   │
   ├── Load Recent Chat (PostgreSQL)
   │
   ├── Build Prompt
   │
   ├── LLM
   │
   ├── Save Assistant Message (PostgreSQL)
   │
   ├── Publish chat.completed (Kafka)
   │
   ▼
Response
   │
   ▼
Visitor


────────────────────────────────────────────────────────
---  background process ----

chat.completed
      │
      ▼
Memory Worker
      │
      ▼
Load User Message
+
Load Assistant Message
(PostgreSQL)
      │
      ▼
Extract Memory (LLM)
      │
      ▼
save == false ?
      │
 ┌────┴────┐
 │         │
Yes        No
 │          │
End         ▼
      Search Existing Memory
      (Milvus)
            │
            ▼
     Memory Consolidation
            │
            ▼
 Update / Insert Memory
(PostgreSQL visitor_knowledge)
            │
            ▼
 Generate Embedding
            │
            ▼
 Update / Insert Vector
(Milvus visitor_knowledge)


──────────────────────────────────
# 1. Visitor Masuk Website
Vue Frontend
      │
      ▼
Check localStorage

Jika belum ada:

visitor_id = ULID()

Simpan ke:

localStorage.setItem("visitor_id", visitorId)

Lalu panggil API:

CreateVisitor

PostgreSQL:

visitors

id = visitor_id
first_seen_at
last_seen_at
total_messages
# 2. Visitor Mengirim Pertanyaan

Contoh:

"Saya pernah cerita project apa saja?"

Frontend mengirim:

{
  "visitor_id": "01K...",
  "session_id": "01K...",
  "message": "Saya pernah cerita project apa saja?"
}

↓

gRPC

↓

ChatService

# 3. Simpan User Message

PostgreSQL

chat_messages

id
session_id
role=user
content=...

# 4. Build Context

ChatService mulai membangun konteks.

## 4.1 Embedding Pertanyaan
Question

↓

Embedding Service

↓

Vector
## 4.2 Search dan_knowledge

Milvus

Collection:
dan_knowledge

Top 5

Misalnya:

Project Portfolio AI

Project Brankazz

Experience Fintech

Kafka

Milvus
## 4.3 Search visitor_knowledge

Milvus

Collection:
visitor_knowledge

Filter:

visitor_id == "01K..."

Top 4

Misalnya:

Visitor sedang membuat Portfolio AI.

Visitor belajar Kafka.

Visitor lebih suka Golang.

Visitor pindah ke Surabaya.
## 4.4 Ambil Chat Terakhir

PostgreSQL

SELECT *
FROM chat_messages
WHERE session_id = ?
ORDER BY created_at DESC
LIMIT 6

Misalnya:

3 user

3 assistant



# 5. Build Prompt

Prompt Builder

SYSTEM

Kamu adalah AI Assistant milik Dani.

--------------------------------

KNOWLEDGE

- Project Portfolio AI
- Project Brankazz
- ...

--------------------------------

VISITOR MEMORY

- Visitor belajar Kafka
- Visitor suka Golang
- ...

--------------------------------

RECENT CHAT

User: ...
Assistant: ...

--------------------------------

QUESTION

Saya pernah cerita project apa saja?


# 6. Generate Jawaban
Prompt

↓

Gemini/OpenAI/Ollama

↓

Response


# 7. Simpan Assistant Message

PostgreSQL

chat_messages

id
session_id
role=assistant
content=...


# 8. Publish Event

Setelah berhasil.

chat.completed

Kafka:

{
  "visitor_id":"01K...",
  "session_id":"01K...",
  "user_message_id":"01K...",
  "assistant_message_id":"01K..."
}


# 9. Return ke Frontend
Response

↓

Vue

↓

Ditampilkan ke visitor.

Background Processing

Ini tidak mengganggu user.



# 10. Memory Worker Consume Event

Kafka

chat.completed

↓

Memory Worker


# 11. Load Conversation

Worker mengambil dari PostgreSQL.

User Message

Assistant Message

Contoh:

User:
Saya sedang belajar Kafka untuk Portfolio AI.

Assistant:
...


# 12. Extract Memory

Prompt khusus:

Extract useful long-term memory.

Ignore:
- greetings
- thanks
- generic questions

Return JSON.

LLM menghasilkan:

{
  "save": true,
  "category": "project",
  "importance": 5,
  "memory": "Visitor sedang mempelajari Kafka untuk membangun Portfolio AI."
}


# 13. Skip Jika Tidak Penting

Misalnya:

Halo

hasil:

{
  "save": false
}

Worker selesai.


# 14. Generate Embedding

Jika save=true

Memory

↓

Embedding

↓

Vector


# 15. Search Existing Memory
Milvus

Collection:
visitor_knowledge

Filter:
visitor_id == "01K..."

Top 3


# 16. Memory Consolidation

Misalnya ditemukan:

Visitor sedang belajar Kafka.

Memory baru:

Visitor sedang mempelajari Kafka untuk membangun Portfolio AI.

LLM menghasilkan:

Visitor sedang membangun Portfolio AI dan mempelajari Kafka sebagai event bus.


# 17. Update Source of Truth

Jika mirip:

UPDATE visitor_knowledge (PostgreSQL)

Jika tidak mirip:

INSERT visitor_knowledge (PostgreSQL)

Contoh isi tabel:

id
visitor_id
category
memory_text
importance
created_at
updated_at


# 18. Generate Embedding

memory_text

↓

Embedding Service

↓

Vector


# 19. Sinkronisasi Milvus

UPSERT

visitor_knowledge_e5
(Milvus)

id
visitor_id
embedding

# Dengan flow ini:

Dengan urutan ini, PostgreSQL selalu menjadi Source of Truth, sedangkan Milvus hanya menjadi vector index. Jika suatu saat kamu mengganti model embedding atau melakukan re-index, cukup membaca ulang seluruh data dari tabel visitor_knowledge di PostgreSQL tanpa kehilangan informasi apa pun.