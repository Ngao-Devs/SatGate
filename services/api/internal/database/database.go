package database

import (
	"database/sql"
	"time"
)

func RunMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS forms (
			id TEXT PRIMARY KEY,
			owner_id TEXT,
			name TEXT NOT NULL,
			domain TEXT NOT NULL,
			amount_sats INTEGER NOT NULL,
			sat_balance INTEGER NOT NULL DEFAULT 0,
			min_amount_sats INTEGER NOT NULL DEFAULT 5,
			created_at DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			form_id TEXT NOT NULL,
			owner_id TEXT,
			sender_name TEXT NOT NULL,
			sender_email TEXT NOT NULL,
			body TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			payment_hash TEXT,
			created_at DATETIME NOT NULL,
			paid_at DATETIME
		);

		CREATE TABLE IF NOT EXISTS invoices (
			id TEXT PRIMARY KEY,
			message_id TEXT NOT NULL,
			amount_sats INTEGER NOT NULL,
			payment_request TEXT NOT NULL,
			payment_hash TEXT,
			status TEXT NOT NULL DEFAULT 'pending',
			preimage TEXT,
			created_at DATETIME NOT NULL,
			paid_at DATETIME,
			expires_at DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS webhook_events (
			id TEXT PRIMARY KEY,
			invoice_id TEXT NOT NULL,
			payload TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	columns := []struct {
		table string
		name  string
		sql   string
	}{
		{"forms", "owner_id", "ALTER TABLE forms ADD COLUMN owner_id TEXT"},
		{"forms", "sat_balance", "ALTER TABLE forms ADD COLUMN sat_balance INTEGER NOT NULL DEFAULT 0"},
		{"forms", "min_amount_sats", "ALTER TABLE forms ADD COLUMN min_amount_sats INTEGER NOT NULL DEFAULT 5"},
		{"messages", "owner_id", "ALTER TABLE messages ADD COLUMN owner_id TEXT"},
	}

	for _, column := range columns {
		if err := ensureColumn(db, column.table, column.name, column.sql); err != nil {
			return err
		}
	}

	if _, err := db.Exec(
		"UPDATE forms SET owner_id = ? WHERE id = ? AND owner_id IS NULL",
		"demo-user",
		"demo-form",
	); err != nil {
		return err
	}

	if _, err := db.Exec(`
		UPDATE messages
		SET owner_id = (
			SELECT owner_id FROM forms WHERE forms.id = messages.form_id
		)
		WHERE owner_id IS NULL
	`); err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_forms_owner_id ON forms(owner_id);
		CREATE INDEX IF NOT EXISTS idx_messages_owner_id ON messages(owner_id);
	`)
	return err
}

func ensureColumn(db *sql.DB, table, column, alterSQL string) error {
	rows, err := db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		return err
	}

	found := false
	for rows.Next() {
		var cid, notNull, primaryKey int
		var name, columnType string
		var defaultValue sql.NullString

		if err := rows.Scan(
			&cid,
			&name,
			&columnType,
			&notNull,
			&defaultValue,
			&primaryKey,
		); err != nil {
			rows.Close()
			return err
		}

		if name == column {
			found = true
		}
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}

	if !found {
		_, err = db.Exec(alterSQL)
	}
	return err
}

func SeedDefaultForm(db *sql.DB) error {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM forms").Scan(&count)
	if count > 0 {
		return nil
	}

	_, err := db.Exec(
		`INSERT INTO forms (
			id, owner_id, name, domain, amount_sats, created_at
		) VALUES (?, ?, ?, ?, ?, ?)`,
		"demo-form", "demo-user", "Sarah's Portfolio", "localhost", 5, time.Now().UTC(),
	)
	return err
}
