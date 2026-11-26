package webhooks

import (
	"context"
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) RecordEvent(
	ctx context.Context,
	eventID string,
	invoiceID string,
	payload string,
	createdAt time.Time,
) (bool, error) {
	result, err := r.db.ExecContext(
		ctx,
		"INSERT OR IGNORE INTO webhook_events (id, invoice_id, payload, created_at) VALUES (?, ?, ?, ?)",
		eventID,
		invoiceID,
		payload,
		createdAt,
	)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (r *Repository) GetPaymentHash(ctx context.Context, invoiceID string) (string, error) {
	var paymentHash string
	err := r.db.QueryRowContext(
		ctx,
		"SELECT payment_hash FROM invoices WHERE id = ?",
		invoiceID,
	).Scan(&paymentHash)

	return paymentHash, err
}
