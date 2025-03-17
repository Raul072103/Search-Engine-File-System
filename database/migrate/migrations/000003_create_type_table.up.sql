CREATE TABLE types (
    id SERIAL PRIMARY KEY,
    file_id BIGINT REFERENCES files(id) ON DELETE CASCADE,
    type INT,
    updated_at TIMESTAMP
);