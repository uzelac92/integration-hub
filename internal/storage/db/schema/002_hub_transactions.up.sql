CREATE TABLE IF NOT EXISTS hub_transactions (
    id BIGSERIAL PRIMARY KEY,
    ref_id TEXT NOT NULL,
    player_id TEXT NOT NULL,
    type TEXT NOT NULL, -- DEBIT or CREDIT
    amount_cents BIGINT NOT NULL,
    currency TEXT NOT NULL,
    operator_status TEXT NOT NULL,
    operator_balance BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS hub_tx_ref_idx ON hub_transactions(ref_id);