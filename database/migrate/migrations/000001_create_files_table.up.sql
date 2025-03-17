CREATE TABLE files (
   id SERIAL PRIMARY KEY,
   path VARCHAR UNIQUE NOT NULL,
   name VARCHAR,
   size BIGINT,
   mode BIGINT,
   extension VARCHAR,
   updated_at TIMESTAMP,
   searchable_tsv TSVECTOR
);