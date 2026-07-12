# Request
Create Profile
Update Profile
Create Project
Update Project
Create Experience
Update Experience
Create Certificate
Update Certificate


# Flow

            Admin (gRPC)
                  │
                  ▼
              Gin API
                  │
                  ▼
      Insert / Update Entity
          (PostgreSQL)
                  │
                  ▼
     Build Knowledge Document
                  │
                  ▼
 Insert / Update knowledge_documents
                  │
                  ▼
      Insert outbox_events
                  │
                  ▼
         COMMIT TRANSACTION
                  │
                  ▼
         Outbox Publisher
                  │
                  ▼
               Kafka
                  │
                  ▼
         Embedding Worker
                  │
                  ▼
 Load Knowledge Document
    from PostgreSQL
                  │
                  ▼
 Knowledge Builder (LLM) Gemini <- AI Agent
                  │
                  ▼
 Structured JSON Response 
{
  "chunks": [
    {
      "title": "...",
      "content": "...",
      "keywords": [...]
    }
  ]
}
                  │
                  ▼
         Golang Processor
                  │
                  ├── Delete old knowledge_chunks
                  ├── Insert new knowledge_chunks
                  ├── Generate embedding
                  └── Update knowledge_documents.status
                  │
                  ▼
          Upsert to Milvus
                  │
                  ▼
 knowledge_documents.status = Embedded


# Tanggung jawab setiap komponen
Komponen	Tanggung Jawab
- Gin API :	CRUD Profile, Project, Experience, Certificate, membangun knowledge_document, membuat outbox_event.

- Outbox Publisher :	Membaca event yang belum dipublish lalu mengirimnya ke Kafka.

- Kafka :	Menjadi event bus/asynchronous queue. Tidak menyimpan data utama.

- Embedding Worker :	Mengambil event, membaca knowledge_document, memproses AI, embedding, dan sinkronisasi Milvus.

- Knowledge Builder (LLM) :	Mengubah satu knowledge_document menjadi beberapa self-contained chunks dalam bentuk JSON.

- Golang Processor :	Menyimpan chunk ke PostgreSQL, menghasilkan embedding, dan melakukan upsert ke Milvus.

- Milvus :	Menyimpan vector embedding untuk semantic search.

Knowledge Builder nantinya modul ini kemungkinan akan berkembang menjadi lebih dari sekadar pemecah dokumen. Misalnya dapat menghasilkan summary, keywords, metadata, atau tags selain chunks, sehingga nama tersebut akan tetap relevan di masa depan.