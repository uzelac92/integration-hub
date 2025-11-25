-- name: GetDueWebhooks :many
SELECT id, event_id, payload, status, attempt_count, next_attempt_at
FROM webhook_outbox
WHERE status = 'PENDING'
  AND next_attempt_at <= NOW()
ORDER BY id
    LIMIT 50;

-- name: MarkWebhookFailed :exec
UPDATE webhook_outbox
SET status = 'FAILED',
    attempt_count = attempt_count + 1,
    next_attempt_at = $2
WHERE id = $1;

-- name: MarkWebhookSuccess :exec
UPDATE webhook_outbox
SET status = 'SUCCESS', attempt_count = attempt_count + 1
WHERE id = $1;