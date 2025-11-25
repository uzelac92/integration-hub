-- name: SaveIdempotency :exec
INSERT INTO idempotency_keys (key, response)
VALUES ($1, $2)
ON CONFLICT (key) DO UPDATE
SET response = EXCLUDED.response;

-- name: GetIdempotency :one
SELECT response
FROM idempotency_keys
WHERE key = $1;

-- name: InsertWebhookOutbox :exec
INSERT INTO webhook_outbox (event_id, payload)
VALUES ($1, $2)
    ON CONFLICT (event_id) DO NOTHING;

-- name: GetWebhookByEventID :one
SELECT id, event_id, payload, status, attempt_count, next_attempt_at, created_at
FROM webhook_outbox
WHERE event_id = $1;

-- name: IncrementWebhookAttempt :exec
UPDATE webhook_outbox
SET attempt_count = attempt_count + 1,
    next_attempt_at = NOW() + INTERVAL '10 seconds'
WHERE id = $1;