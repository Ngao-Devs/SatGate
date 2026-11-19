package database

import (
	"database/sql"
	"time"
)

func RunMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS forms (
			id          TEXT PRIMARY KEY,
			owner_id    TEXT,
			name        TEXT NOT NULL,
			domain      TEXT NOT NULL,
			amount_sats INTEGER NOT NULL,
			created_at  DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS messages (
			id           TEXT PRIMARY KEY,
			form_id      TEXT NOT NULL,
			sender_name  TEXT NOT NULL,
			sender_email TEXT NOT NULL,
			body         TEXT NOT NULL,
			status       TEXT NOT NULL DEFAULT 'pending',
			payment_hash TEXT,
			created_at   DATETIME NOT NULL,
			paid_at      DATETIME
		);

		CREATE TABLE IF NOT EXISTS invoices (
			id              TEXT PRIMARY KEY,
			message_id      TEXT NOT NULL,
			amount_sats     INTEGER NOT NULL,
			payment_request TEXT NOT NULL,
			payment_hash    TEXT,
			status          TEXT NOT NULL DEFAULT 'pending',
			preimage        TEXT,
			created_at      DATETIME NOT NULL,
			paid_at         DATETIME,
			expires_at      DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS webhook_events (
			id         TEXT PRIMARY KEY,
			invoice_id TEXT NOT NULL,
			payload    TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	// Upgrade existing databases without destroying existing data.
	rows, err := db.Query("PRAGMA table_info(forms)")
	if err != nil {
		return err
	}

	hasOwnerID := false

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

		if name == "owner_id" {
			hasOwnerID = true
		}
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	if !hasOwnerID {
		if _, err := db.Exec("ALTER TABLE forms ADD COLUMN owner_id TEXT"); err != nil {
			return err
		}
	}

	// Preserve ownership of the existing development demo form.
	if _, err := db.Exec(
		"UPDATE forms SET owner_id = ? WHERE id = ? AND owner_id IS NULL",
		"demo-user",
		"demo-form",
	); err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_forms_owner_id
		ON forms(owner_id)
	`)

	return err
}

func SeedDefaultForm(db *sql.DB) error {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM forms").Scan(&count)
	if count > 0 {
		return nil
	}

	_, err := db.Exec(
		"INSERT INTO forms (id, owner_id, name, domain, amount_sats, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		"demo-form", "demo-user", "Sarah's Portfolio", "localhost", 5, time.Now().UTC(),
	)
	return err
}
