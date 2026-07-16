saya akan membaginya menjadi dua jalur:

Online Path → menjawab user secepat mungkin.
Background Path → membangun memori jangka panjang.

Dengan begitu latency chat tetap rendah.

# Online Path (Realtime)

                User
                  │
                  ▼
              Gin API
                  │
                  ▼
            Create Session
            Create Message
          (PostgreSQL)
                  │
                  ▼
      Embed User Question
                  │
      ┌───────────┴────────────┐
      ▼                        ▼

Knowledge Search Visitor Memory Search
(Milvus) (Milvus)
top 5 knowledge top 4 visitor memory
│ │
└───────────┬────────────┘
▼
Build Prompt
│
▼
LLM
│
▼
Save Assistant Message
(PostgreSQL)
│
▼
Return Response

User mendapat jawaban secepat mungkin.

# Background Path

Setelah AI selesai menjawab

Save Assistant Message
│
▼
Publish Kafka Event
chat.completed

Worker menerima event.

Kafka
│
▼
Memory Worker

Worker mengambil:

User Message
Assistant Message
Visitor ID
Session ID

Lalu meminta AI membuat memory.

Prompt Memory

Misalnya prompt:

You are a memory extractor.

Extract only information useful for future conversations.

Ignore:

- greetings
- thanks
- jokes
- generic questions

Return JSON:

{
"save": true|false,
"importance":1-5,
"category":"",
"memory":""
}

Maximum 250 characters.

Contoh

User

Saya sedang membangun Portfolio AI memakai Kafka.

AI

...

Output

{
"save": true,
"importance": 5,
"category": "project",
"memory": "Visitor sedang membangun Portfolio AI menggunakan Kafka."
}

Kalau

User

Makasih.

AI

Sama-sama.

Output

{
"save": false
}

Worker berhenti.

Kalau save=true

Worker membuat embedding.

Memory Text
│
▼
Embedding Model
│
▼
768 Vector

Kemudian insert ke Milvus.

visitor_knowledge

id
visitor_id
session_id
category
importance
memory
embedding
created_at
Tetapi saya akan menambah satu langkah lagi

Sebelum insert

Cari memory yang mirip.

Memory Baru
│
▼
Milvus Search

filter:

visitor_id == xxxx

top 3

Misalnya hasilnya

Visitor sedang belajar Kafka.

score = 0.93

Padahal memory baru

Visitor sedang membangun Portfolio AI menggunakan Kafka.

Sangat mirip.

Daripada insert lagi

Worker melakukan merge.

Prompt

Combine these memories into one.

Memory A:
Visitor sedang belajar Kafka.

Memory B:
Visitor sedang membangun Portfolio AI menggunakan Kafka.

Return one memory.

Hasil

Visitor sedang membangun Portfolio AI dan mempelajari Kafka sebagai event bus.

Update Milvus.

Tidak insert baru.

Kalau similarity kecil

misalnya

Visitor tinggal di Surabaya.

Tidak ada yang mirip.

Insert baru.

Retrieval

Ketika user bertanya

User Question
│
▼
Embedding
│
├──────────────┐
▼ ▼
Knowledge Visitor Memory
Top 5 Top 4

Prompt

Knowledge

- ...
- ...

Visitor Memory

- sedang membangun Portfolio AI
- tinggal di Surabaya
- lebih suka Golang

Question

...
