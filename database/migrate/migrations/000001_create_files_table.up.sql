CREATE TABLE files (
   id SERIAL PRIMARY KEY,
   path VARCHAR UNIQUE NOT NULL,
   name VARCHAR,
   size BIGINT,
   mode BIGINT,
   extension VARCHAR,
   file_id BIGINT UNIQUE,
   parent_id BIGINT,
   rank DOUBLE PRECISION,
   hash CHAR(32),
   updated_at TIMESTAMP
);