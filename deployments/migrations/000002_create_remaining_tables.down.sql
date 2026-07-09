-- deployments/migrations/000002_create_remaining_tables.down.sql

-- Drop tables in reverse order of creation (dependencies first)
DROP TABLE IF EXISTS retrieval_logs;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS outbox_events;
DROP TABLE IF EXISTS ai_tools;
DROP TABLE IF EXISTS prompts;
DROP TABLE IF EXISTS ai_models;
DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS chat_sessions;
DROP TABLE IF EXISTS visitor_memories;
DROP TABLE IF EXISTS visitors;
DROP TABLE IF EXISTS knowledge_chunks;
DROP TABLE IF EXISTS knowledge_documents;
DROP TABLE IF EXISTS blogs;
DROP TABLE IF EXISTS skills;
DROP TABLE IF EXISTS certificates;
DROP TABLE IF EXISTS project_technologies;
DROP TABLE IF EXISTS technologies;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS experiences;

-- Revert columns added to profiles table
ALTER TABLE profiles DROP COLUMN IF EXISTS availability;
ALTER TABLE profiles DROP COLUMN IF EXISTS timezone;
