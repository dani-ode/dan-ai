1. profiles
id #ULID
full_name
headline
bio
email
phone
location
github
linkedin
website
avatar
resume_url
availability #Available, Busy, Not Looking
timezone
created_at
updated_at

2. experiences
id #ULID
company
position
employment_type
start_date
end_date
current_job
location
description
display_order
company_logo
skills #JSON
remote_type #Remote, Hybrid, Onsite
created_at
updated_at

3. projects
id #ULID
slug
title
summary
description
architecture
repository_url
demo_url
thumbnail
featured
status # Draft, Published, Archived
github_stars
github_last_commit
read_time
created_at
updated_at

4. technologies
id #ULID
name
category
icon
color
official_url
logo
created_at

5. project_technologies #Many-to-many.
project_id
technology_id
display_order

6. certificates
id #ULID
title
issuer
issue_date
expiration_date
credential_id
credential_url
thumbnail
skills #JSON
issuer_logo
created_at

7. skills
id #ULID
technology_id
display_order
level
years
favorite
created_at


8. knowledge_documents #Ini inti AI
id #ULID
source_type #profile, experience, project, certificate, blog, manual
source_id
title
content
checksum
version
status # Pending, Embedding, Embedded, Failed
embedding_model
last_embedded_at
created_at
updated_at

9. knowledge_chunks
id #ULID
document_id
chunk_index
content
token_count
embedding_model
created_at

10. visitors
id #ULID
first_seen_at
last_seen_at
total_messages
created_at
updated_at

11. visitor_memories
id #ULID
visitor_id
key
value
created_at
updated_at

12. chat_sessions
id #ULID
visitor_id
prompt_id
title
started_at
ended_at
created_at
updated_at

13. chat_messages
id #ULID
session_id
role #system, user, assistant, tool
content
model
prompt_tokens
completion_tokens
latency_ms
status #Pending, Streaming, Completed, Error
created_at
updated_at

14. prompts
id #ULID
name
system_prompt
description
model_id
active
version
created_at

15. ai_models
id #ULID
name
provider
temperature
max_tokens
context_window
supports_tools
supports_stream
enabled

16. ai_tools #Kalau nanti Agent punya tool.
id #ULID
name
tool_type
config #JSON
description
enabled

17. outbox_events #Kafka Event.
id #ULID
aggregate
aggregate_id
event_type
payload
published
retry_count
failed_reason
published_at
created_at