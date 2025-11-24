CREATE TABLE IF NOT EXISTS idempotency_keys (
    id SERIAL PRIMARY KEY,
    key TEXT UNIQUE NOT NULL,
    response JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_idempotency_key ON idempotency_keys(key);