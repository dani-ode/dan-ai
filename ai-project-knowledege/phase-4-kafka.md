Kafka

Saya lebih suka topic berdasarkan business event.

Misalnya

portfolio.profile.updated

portfolio.project.created

portfolio.project.updated

portfolio.experience.created

portfolio.experience.updated

portfolio.certificate.created

portfolio.certificate.updated

Tetapi...

Menurut saya ini terlalu banyak.

Saya lebih suka

Topic
portfolio.knowledge

Semua event knowledge masuk sini.

Payload

{
    "aggregate":"project",
    "aggregate_id":"01K...",
    "event_type":"updated"
}

atau

{
    "aggregate":"profile",
    "aggregate_id":"01K...",
    "event_type":"updated"
}

Worker tinggal switch

switch aggregate {
case "project":
...
case "profile":
...
}

Jauh lebih sederhana.

Kemudian

Topic lain

portfolio.events

untuk Event Worker.

Misalnya

GitHub Sync

Analytics

Discord Notification

Email

Thumbnail

dll

Jadi total cukup

portfolio.knowledge

portfolio.events
Consumer Group

Embedding Worker

portfolio-embedding-worker

Event Worker

portfolio-event-worker

Nanti kalau worker embedding ada 3 container

Embedding Worker 1

Embedding Worker 2

Embedding Worker 3

mereka tetap memakai

portfolio-embedding-worker

Kafka otomatis membagi partition.

Partition

Saat development

Saya sarankan

Partition = 1

Sudah cukup.

Karena

hanya satu VPS
hanya satu worker
traffic kecil

Kalau nanti production

baru

portfolio.knowledge

Partition = 3

atau

Partition = 6

Worker

Worker 1

Worker 2

Worker 3

Kafka akan load balancing sendiri.

Jadi untuk sekarang
Kafka
Topic

portfolio.knowledge
portfolio.events

Consumer Group

portfolio-embedding-worker
portfolio-event-worker

Partition

1
