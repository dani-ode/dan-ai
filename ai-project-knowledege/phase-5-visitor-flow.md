                    USER
                      │
                      ▼
                 Gin API
                      │
                      ▼
           Save User Message (Postgres)
                      │
                      ▼
          Embed Current Question
                      │
        ┌─────────────┴─────────────┐
        ▼                           ▼

Knowledge Search Visitor Memory Search
(Milvus) (Milvus)
│ │
└─────────────┬─────────────┘
▼
Build Prompt
│
▼
LLM
│
▼
Save Assistant Message (Postgres)
│
▼
Publish chat.completed
(Kafka)
│
──────────────────────┼──────────────────────────────
│
Memory Worker
│
▼
Extract Memory (LLM)
│
save == false ?
┌───────┴────────┐
│ │
Yes No
│ ▼
Finish Embed Memory
│
▼
Search Similar Memory
│
similarity > threshold ?
┌────────┴────────┐
│ │
Yes No
│ │
Merge & Update Insert New
│ │
└────────┬────────┘
▼
visitor_knowledge
