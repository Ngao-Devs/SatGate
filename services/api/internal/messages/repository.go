package messages

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListPaidByOwner(ctx context.Context, ownerID string) ([]Message, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT m.id, m.form_id, m.sender_name, m.sender_email, m.body, m.status,
		       COALESCE(i.amount_sats, 0), m.created_at, m.paid_at
		FROM messages m
		JOIN forms f ON f.id = m.form_id
		LEFT JOIN invoices i ON i.message_id = m.id
		WHERE m.status = 'paid' AND m.owner_id = ?
		ORDER BY m.paid_at DESC, m.created_at DESC
	`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []Message{}

	for rows.Next() {
		var message Message
		if err := rows.Scan(
			&message.ID,
			&message.FormID,
			&message.SenderName,
			&message.SenderEmail,
			&message.Body,
			&message.Status,
			&message.AmountSats,
			&message.CreatedAt,
			&message.PaidAt,
		); err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
