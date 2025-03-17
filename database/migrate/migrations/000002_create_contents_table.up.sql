CREATE TABLE contents (
  id SERIAL PRIMARY KEY,
  file_id BIGINT REFERENCES files(id) ON DELETE CASCADE,
  content_text TEXT,
  content_bytes BYTEA,
  searchable_tsv TSVECTOR,
  updated_at TIMESTAMP
);