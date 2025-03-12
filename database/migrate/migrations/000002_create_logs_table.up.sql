CREATE TABLE IF NOT EXISTS logs (
    id bigserial PRIMARY KEY,
    level varchar(10),
    msg text,
    logger varchar(255),
    caller varchar(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    context jsonb
);