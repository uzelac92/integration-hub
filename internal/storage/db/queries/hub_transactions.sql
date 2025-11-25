-- name: InsertHubTransaction :exec
INSERT INTO hub_transactions (
    ref_id,
    player_id,
    type,
    amount_cents,
    currency,
    operator_status,
    operator_balance
)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: ListHubTransactions :many
SELECT *
FROM hub_transactions
ORDER BY id DESC;