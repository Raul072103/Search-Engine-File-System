CREATE TABLE IF NOT EXISTS files (
    id bigserial PRIMARY KEY,
    path text UNIQUE NOT NULL,
    name text NOT NULL,
    size bigint,
    is_dir bool,
    mode integer,
    extension text,
    updated_at TIMESTAMP NOT NULL,
    content text,
    searchable_tsv tsvector
);