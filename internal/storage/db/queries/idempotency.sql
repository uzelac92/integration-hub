-- name: SaveIdempotency :exec
INSERT INTO idempotency_keys (key, response)
VALUES ($1, $2)
ON CONFLICT (key) DO UPDATE
SET response = EXCLUDED.response;

-- name: GetIdempotency :one
SELECT response
FROM idempotency_keys
WHERE key = $1;