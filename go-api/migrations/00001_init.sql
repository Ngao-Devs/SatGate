-- +goose Up
CREATE TABLE IF NOT EXISTS sites (
    id VARCHAR(255) PRIMARY KEY,
    owner_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) NOT NULL DEFAULT '',
    sat_balance BIGINT NOT NULL DEFAULT 0,
    min_amount_sats BIGINT NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS invoices (
    payment_hash VARCHAR(255) PRIMARY KEY,
    payment_request TEXT NOT NULL,
    site_id VARCHAR(255) NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    amount_sats BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS submissions (
    id VARCHAR(255) PRIMARY KEY,
    site_id VARCHAR(255) NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'verified',
    paid_sats BIGINT NOT NULL,
    timestamp VARCHAR(255) NOT NULL,
    payment_hash VARCHAR(255) NOT NULL UNIQUE,
    preimage VARCHAR(255) NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_submissions_site_id ON submissions(site_id);
CREATE INDEX IF NOT EXISTS idx_sites_owner_id ON sites(owner_id);

-- +goose Down
DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS sites;