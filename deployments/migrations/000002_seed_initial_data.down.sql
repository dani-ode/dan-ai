-- deployments/migrations/000002_seed_initial_data.down.sql

DELETE FROM embedding_profiles WHERE name IN ('e5', 'openai-small');
DELETE FROM prompts WHERE name IN ('Knowledge Chunker', 'Default Assistant Prompt', 'Memory Extractor', 'Memory Consolidator');
DELETE FROM ai_models WHERE name IN ('gemini-3.1-flash-lite', 'gemini-2.0-flash-lite', 'gemini-embedding-2', 'gpt-4o-mini', 'text-embedding-3-small', 'text-embedding-3-large');
