Milvus

saya sudah punya milvus yang jalan dengan database "portfolio" tanpa auth di port "127.0.0.1:19530"
40c0af03f4fc   milvusdb/milvus:2.4-latest   "/tini -- milvus run…"   9 months ago   Up 38 minutes          0.0.0.0:9091->9091/tcp, [::]:9091->9091/tcp, 0.0.0.0:19530->19530/tcp, [::]:19530->19530/tcp   milvus-standalone
1a9a7d3d955d   quay.io/coreos/etcd:v3.5.5   "etcd -advertise-cli…"   9 months ago   Up 39 minutes          0.0.0.0:2379->2379/tcp, [::]:2379->2379/tcp                                                    milvus-etcd
c1c66e0e576f   minio/minio                  "/usr/bin/docker-ent…"   9 months ago   Up 39 minutes          0.0.0.0:9000->9000/tcp, [::]:9000->9000/tcp                                                    milvus-minio



buat satu collection saja
portfolio_knowledge

Karena semua yang dicari AI adalah Knowledge.

Nanti filtering menggunakan metadata.

Schema
Collection

portfolio_knowledge

Field

id                 VARCHAR(26)     Primary Key
document_id        VARCHAR(26)
chunk_id           VARCHAR(26)

source_type        VARCHAR(30)
source_id          VARCHAR(26)

embedding          FLOAT_VECTOR

title              VARCHAR(200)

content            VARCHAR(4000)

keywords           JSON

created_at         INT64

Kalau memakai Milvus terbaru, metadata JSON juga bisa.

Misalnya

{
    "project_id":"...",
    "technology":["Go","Kafka"],
    "featured":true
}

Tetapi menurut saya cukup

source_type
source_id
document_id
chunk_id

sisanya tetap di PostgreSQL.

Contoh isi
id

01K.....

document_id

01K.....

chunk_id

01K.....

source_type

project

source_id

01K.....

title

Portfolio AI

content

Portfolio AI menggunakan Gin,
Kafka,
Milvus,
Gemini...

embedding

[]

Index

Misalnya

Metric

COSINE

karena embedding Gemini/OpenAI umumnya memakai cosine similarity.

Search

Misalnya user bertanya

"Bagaimana arsitektur backendmu?"

Worker

↓

Embedding

↓

Milvus

↓

TopK = 5

↓

Return

chunk 4

chunk 7

chunk 1

chunk 9

chunk 15


# Milvus Collection

portfolio_knowledge

Primary Key

id

Vector

embedding

Metadata

document_id
chunk_id
source_type
source_id
title
content
created_at

Index

COSINE
Satu penyederhanaan lagi yang saya rekomendasikan

Karena knowledge_documents sudah menyimpan source_type dan source_id, saya tidak akan menduplikasi terlalu banyak metadata di Milvus. Saya cukup menyimpan apa yang dibutuhkan untuk pencarian dan identifikasi hasil.

Collection portfolio_knowledge:

Field	Keterangan
chunk_id (Primary Key)	ID unik chunk, sama dengan knowledge_chunks.id.
document_id	Relasi ke knowledge_documents.
source_type	profile, project, experience, certificate, dll.
source_id	ID entity asal (projects.id, profiles.id, dll.).
embedding	Vector embedding.

Ketika Milvus mengembalikan chunk_id, backend tinggal melakukan:

Milvus
    │
    ▼
chunk_id
    │
    ▼
SELECT * FROM knowledge_chunks
WHERE id = ?

Dengan cara ini, Milvus hanya berfungsi sebagai vector index, sedangkan seluruh isi teks, metadata lengkap, dan source of truth tetap berada di PostgreSQL. Pendekatan ini membuat sinkronisasi lebih sederhana dan menghindari duplikasi data yang tidak perlu.