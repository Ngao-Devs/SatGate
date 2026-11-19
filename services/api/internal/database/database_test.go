package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	db.SetMaxOpenConns(1)

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestRunMigrationsCreatesTables(t *testing.T) {
	db := openTestDB(t)

	if err := RunMigrations(db); err != nil {
		t.Fatalf("RunMigrations returned error: %v", err)
	}

	for _, table := range []string{"forms", "messages", "invoices", "webhook_events"} {
		var name string
		err := db.QueryRow(
			"SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?",
			table,
		).Scan(&name)

		if err != nil {
			t.Errorf("expected table %q to exist: %v", table, err)
		}
	}
}

func TestRunMigrationsUpgradesExistingFormsTable(t *testing.T) {
	db := openTestDB(t)

	_, err := db.Exec(`
		CREATE TABLE forms (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			domain TEXT NOT NULL,
			amount_sats INTEGER NOT NULL,
			created_at DATETIME NOT NULL
		);

		INSERT INTO forms (id, name, domain, amount_sats, created_at)
		VALUES ('demo-form', 'Demo', 'localhost', 5, CURRENT_TIMESTAMP);
	`)
	if err != nil {
		t.Fatalf("create legacy forms table: %v", err)
	}

	if err := RunMigrations(db); err != nil {
		t.Fatalf("RunMigrations returned error: %v", err)
	}

	var ownerID string
	if err := db.QueryRow(
		"SELECT owner_id FROM forms WHERE id = 'demo-form'",
	).Scan(&ownerID); err != nil {
		t.Fatalf("read migrated owner: %v", err)
	}

	if ownerID != "demo-user" {
		t.Errorf("expected demo-user owner, got %q", ownerID)
	}

	var indexName string
	if err := db.QueryRow(
		"SELECT name FROM sqlite_master WHERE type = 'index' AND name = 'idx_forms_owner_id'",
	).Scan(&indexName); err != nil {
		t.Errorf("expected owner index to exist: %v", err)
	}
}

func TestRunMigrationsPreservesExistingOwner(t *testing.T) {
	db := openTestDB(t)

	if err := RunMigrations(db); err != nil {
		t.Fatalf("initial migration: %v", err)
	}

	_, err := db.Exec(`
		INSERT INTO forms (id, owner_id, name, domain, amount_sats, created_at)
		VALUES ('existing-form', 'existing-user', 'Existing', 'example.com', 10, CURRENT_TIMESTAMP)
	`)
	if err != nil {
		t.Fatalf("insert existing form: %v", err)
	}

	if err := RunMigrations(db); err != nil {
		t.Fatalf("re-run migrations: %v", err)
	}

	var ownerID string
	if err := db.QueryRow(
		"SELECT owner_id FROM forms WHERE id = 'existing-form'",
	).Scan(&ownerID); err != nil {
		t.Fatalf("read owner: %v", err)
	}

	if ownerID != "existing-user" {
		t.Errorf("expected existing-user owner, got %q", ownerID)
	}
}

func TestSeedDefaultForm(t *testing.T) {
	db := openTestDB(t)

	if err := RunMigrations(db); err != nil {
		t.Fatalf("RunMigrations: %v", err)
	}

	if err := SeedDefaultForm(db); err != nil {
		t.Fatalf("SeedDefaultForm returned error: %v", err)
	}

	if err := SeedDefaultForm(db); err != nil {
		t.Fatalf("second SeedDefaultForm returned error: %v", err)
	}

	var count int
	if err := db.QueryRow(
		"SELECT COUNT(*) FROM forms WHERE id = 'demo-form'",
	).Scan(&count); err != nil {
		t.Fatalf("count seeded forms: %v", err)
	}

	if count != 1 {
		t.Errorf("expected one demo form, got %d", count)
	}
}
