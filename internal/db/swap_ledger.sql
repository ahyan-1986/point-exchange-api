CREATE TABLE swap_ledger (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_partner_id VARCHAR(64) NOT NULL,
    source_external_id VARCHAR(128) NOT NULL,
    source_customer_id VARCHAR(64) NOT NULL,
    source_points FLOAT NOT NULL,
    usd_value FLOAT NOT NULL,
    target_partner_id VARCHAR(64) NOT NULL,
    target_customer_id VARCHAR(64) NOT NULL,
    target_points FLOAT NOT NULL,
    status VARCHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_swap_ledger_source_external_id ON swap_ledger(source_external_id);
CREATE INDEX idx_swap_ledger_status ON swap_ledger(status);
