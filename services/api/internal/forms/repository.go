package forms

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("form not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListByOwner(ctx context.Context, ownerID string) ([]Form, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, domain, amount_sats, sat_balance, min_amount_sats, created_at
		FROM forms
		WHERE owner_id = ?
		ORDER BY created_at DESC
	`, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list forms: %w", err)
	}
	defer rows.Close()

	forms := make([]Form, 0)

	for rows.Next() {
		var form Form
		if err := rows.Scan(
			&form.ID,
			&form.Name,
			&form.Domain,
			&form.AmountSats,
			&form.SatBalance,
			&form.MinAmountSats,
			&form.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan form: %w", err)
		}

		forms = append(forms, form)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate forms: %w", err)
	}

	return forms, nil
}

func (r *Repository) Create(ctx context.Context, ownerID string, form Form) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO forms (
			id, owner_id, name, domain, amount_sats,
			sat_balance, min_amount_sats, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		form.ID,
		ownerID,
		form.Name,
		form.Domain,
		form.AmountSats,
		form.SatBalance,
		form.MinAmountSats,
		form.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create form: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Form, error) {
	var form Form

	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, domain, amount_sats, sat_balance, min_amount_sats, created_at
		FROM forms
		WHERE id = ?
	`, id).Scan(
		&form.ID,
		&form.Name,
		&form.Domain,
		&form.AmountSats,
		&form.SatBalance,
		&form.MinAmountSats,
		&form.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get form: %w", err)
	}

	return &form, nil
}

func (r *Repository) GetByIDAndOwner(ctx context.Context, id, ownerID string) (*Form, error) {
	var form Form

	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, domain, amount_sats, sat_balance, min_amount_sats, created_at
		FROM forms
		WHERE id = ? AND owner_id = ?
	`, id, ownerID).Scan(
		&form.ID,
		&form.Name,
		&form.Domain,
		&form.AmountSats,
		&form.SatBalance,
		&form.MinAmountSats,
		&form.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get owned form: %w", err)
	}

	return &form, nil
}

func (r *Repository) Update(ctx context.Context, ownerID string, form Form) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE forms
		SET name = ?, domain = ?, amount_sats = ?, min_amount_sats = ?
		WHERE id = ? AND owner_id = ?
	`,
		form.Name,
		form.Domain,
		form.AmountSats,
		form.MinAmountSats,
		form.ID,
		ownerID,
	)
	if err != nil {
		return fmt.Errorf("update form: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check updated form: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id, ownerID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin form deletion: %w", err)
	}
	defer tx.Rollback()

	queries := []string{
		`DELETE FROM webhook_events
		 WHERE invoice_id IN (
			SELECT id FROM invoices
			WHERE message_id IN (
				SELECT id FROM messages
				WHERE form_id IN (
					SELECT id FROM forms WHERE id = ? AND owner_id = ?
				)
			)
		 )`,
		`DELETE FROM invoices
		 WHERE message_id IN (
			SELECT id FROM messages
			WHERE form_id IN (
				SELECT id FROM forms WHERE id = ? AND owner_id = ?
			)
		 )`,
		`DELETE FROM messages
		 WHERE form_id IN (
			SELECT id FROM forms WHERE id = ? AND owner_id = ?
		 )`,
	}

	for _, query := range queries {
		if _, err := tx.ExecContext(ctx, query, id, ownerID); err != nil {
			return fmt.Errorf("delete related form records: %w", err)
		}
	}

	result, err := tx.ExecContext(ctx, `
		DELETE FROM forms
		WHERE id = ? AND owner_id = ?
	`, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete form: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check deleted form: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit form deletion: %w", err)
	}

	return nil
}
