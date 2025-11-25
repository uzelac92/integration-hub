-- name: InsertHubWebhookEvent :exec
INSERT INTO hub_webhook_events (event_id, payload)
VALUES ($1, $2)
    ON CONFLICT (event_id) DO NOTHING;

-- name: HubWebhookEventExists :one
SELECT EXISTS (
    SELECT 1
    FROM hub_webhook_events
    WHERE event_id = $1
);