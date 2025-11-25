-- name: GetDueWebhooks :many
SELECT id, event_id, payload, status, attempt_count, next_attempt_at
FROM webhook_outbox
WHERE status = 'PENDING'
  AND next_attempt_at <= NOW()
ORDER BY id
    LIMIT 50;

-- name: MarkWebhookSent :exec
UPDATE webhook_outbox
SET status = 'SENT'
WHERE id = $1;

-- name: MarkWebhookFailed :exec
UPDATE webhook_outbox
SET attempt_count = attempt_count + 1,
    next_attempt_at = NOW() + INTERVAL '10 seconds'
WHERE id = $1;