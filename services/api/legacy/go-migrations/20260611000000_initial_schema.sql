-- Initial schema for SatGate

CREATE TABLE forms (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    domain TEXT NOT NULL,
    amount_sats INTEGER NOT NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    form_id TEXT NOT NULL REFERENCES forms(id),
    sender_name TEXT NOT NULL,
    sender_email TEXT NOT NULL,
    body TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TEXT NOT NULL,
    paid_at TEXT
);

CREATE TABLE invoices (
    id TEXT PRIMARY KEY,
    message_id TEXT NOT NULL REFERENCES messages(id),
    amount_sats INTEGER NOT NULL,
    payment_request TEXT NOT NULL,
    payment_hash TEXT NOT NULL,
    status TEXT NOT NULL,
    preimage TEXT,
    created_at TEXT NOT NULL,
    paid_at TEXT,
    expires_at TEXT NOT NULL
);

CREATE TABLE webhook_events (
    id TEXT PRIMARY KEY,
    invoice_id TEXT NOT NULL,
    payload TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE INDEX idx_messages_status ON messages(status);
CREATE INDEX idx_messages_form_id ON messages(form_id);
CREATE INDEX idx_invoices_message_id ON invoices(message_id);
CREATE INDEX idx_invoices_status ON invoices(status);
