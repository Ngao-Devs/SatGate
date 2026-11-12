-- Additional indices for performance

CREATE INDEX idx_messages_paid_at ON messages(paid_at);
CREATE INDEX idx_forms_created_at ON forms(created_at);
