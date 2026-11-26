package invoices

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetStatus(ctx context.Context, id string) (*InvoiceStatusResponse, error) {
	var response InvoiceStatusResponse

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, message_id, status, paid_at FROM invoices WHERE id = ?",
		id,
	).Scan(
		&response.InvoiceID,
		&response.MessageID,
		&response.Status,
		&response.PaidAt,
	)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (r *Repository) MarkPaid(ctx context.Context, invoiceID string, preimage *string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	paidAt := time.Now().UTC()

	result, err := tx.Exec(`
		UPDATE invoices
		SET status = 'paid', paid_at = ?, preimage = COALESCE(?, preimage)
		WHERE id = ? AND status != 'paid'
	`, paidAt, preimage, invoiceID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		var exists int
		err = tx.QueryRow(
			"SELECT 1 FROM invoices WHERE id = ?",
			invoiceID,
		).Scan(&exists)
		if err != nil {
			return fmt.Errorf("invoice not found: %w", err)
		}
		return tx.Commit()
	}

	var amountSats int64
	var messageID, formID string

	err = tx.QueryRow(`
		SELECT i.amount_sats, m.id, m.form_id
		FROM invoices i
		JOIN messages m ON m.id = i.message_id
		WHERE i.id = ?
	`, invoiceID).Scan(&amountSats, &messageID, &formID)
	if err != nil {
		return fmt.Errorf("invoice or message not found: %w", err)
	}

	_, err = tx.Exec(
		"UPDATE messages SET status = 'paid', paid_at = ? WHERE id = ?",
		paidAt,
		messageID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE forms SET sat_balance = sat_balance + ? WHERE id = ?",
		amountSats,
		formID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
