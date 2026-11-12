-- Add payment_hash to messages table for better tracking and integrity

ALTER TABLE messages ADD COLUMN payment_hash TEXT;
CREATE UNIQUE INDEX idx_messages_payment_hash ON messages(payment_hash);
