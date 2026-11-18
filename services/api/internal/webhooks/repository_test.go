package webhooks

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func newTestWebhookRepository(t *testing.T) (*sql.DB, *Repository) {
	t.Helper()

	db, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "webhooks.db"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE invoices (
			id TEXT PRIMARY KEY,
			payment_hash TEXT
		);

		CREATE TABLE webhook_events (
			id TEXT PRIMARY KEY,
			invoice_id TEXT NOT NULL,
			payload TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);

		INSERT INTO invoices (id, payment_hash)
		VALUES ('invoice-1', 'hash-123');
	`)
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db, NewRepository(db)
}

func TestRecordEvent(t *testing.T) {
	db, repository := newTestWebhookRepository(t)

	recorded, err := repository.RecordEvent(
		context.Background(),
		"event-1",
		"invoice-1",
		`{"status":"paid"}`,
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !recorded {
		t.Fatal("expected event to be recorded")
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM webhook_events").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 event, got %d", count)
	}
}

func TestRecordDuplicateEvent(t *testing.T) {
	_, repository := newTestWebhookRepository(t)

	ctx := context.Background()
	createdAt := time.Now().UTC()

	if recorded, err := repository.RecordEvent(ctx, "event-1", "invoice-1", `{}`, createdAt); err != nil || !recorded {
		t.Fatalf("first event: recorded=%t err=%v", recorded, err)
	}

	recorded, err := repository.RecordEvent(ctx, "event-1", "invoice-1", `{}`, createdAt)
	if err != nil {
		t.Fatal(err)
	}
	if recorded {
		t.Fatal("expected duplicate event not to be recorded")
	}
}

func TestGetPaymentHash(t *testing.T) {
	_, repository := newTestWebhookRepository(t)

	hash, err := repository.GetPaymentHash(context.Background(), "invoice-1")
	if err != nil {
		t.Fatal(err)
	}
	if hash != "hash-123" {
		t.Fatalf("expected hash-123, got %s", hash)
	}
}

func TestGetPaymentHashMissingInvoice(t *testing.T) {
	_, repository := newTestWebhookRepository(t)

	_, err := repository.GetPaymentHash(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error for missing invoice")
	}
}
