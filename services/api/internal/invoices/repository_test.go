package invoices

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func newTestInvoiceRepository(t *testing.T) (*sql.DB, *Repository) {
	t.Helper()

	db, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "invoices.db"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE messages (
			id TEXT PRIMARY KEY,
			status TEXT NOT NULL,
			paid_at DATETIME
		);

		CREATE TABLE invoices (
			id TEXT PRIMARY KEY,
			message_id TEXT NOT NULL,
			status TEXT NOT NULL,
			preimage TEXT,
			paid_at DATETIME
		);

		INSERT INTO messages (id, status) VALUES ('message-1', 'pending');
		INSERT INTO invoices (id, message_id, status)
		VALUES ('invoice-1', 'message-1', 'pending');
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

func TestGetStatus(t *testing.T) {
	_, repository := newTestInvoiceRepository(t)

	response, err := repository.GetStatus(context.Background(), "invoice-1")
	if err != nil {
		t.Fatal(err)
	}

	if response.InvoiceID != "invoice-1" {
		t.Fatalf("expected invoice-1, got %s", response.InvoiceID)
	}

	if response.MessageID != "message-1" {
		t.Fatalf("expected message-1, got %s", response.MessageID)
	}

	if response.Status != "pending" {
		t.Fatalf("expected pending, got %s", response.Status)
	}
}

func TestGetStatusMissingInvoice(t *testing.T) {
	_, repository := newTestInvoiceRepository(t)

	_, err := repository.GetStatus(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error for missing invoice")
	}
}

func TestMarkPaidUpdatesInvoiceAndMessage(t *testing.T) {
	db, repository := newTestInvoiceRepository(t)

	preimage := "test-preimage"
	if err := repository.MarkPaid(context.Background(), "invoice-1", &preimage); err != nil {
		t.Fatal(err)
	}

	var invoiceStatus, storedPreimage, messageStatus string
	var invoicePaidAt, messagePaidAt time.Time

	err := db.QueryRow(`
		SELECT i.status, i.preimage, i.paid_at, m.status, m.paid_at
		FROM invoices i
		JOIN messages m ON m.id = i.message_id
		WHERE i.id = 'invoice-1'
	`).Scan(
		&invoiceStatus,
		&storedPreimage,
		&invoicePaidAt,
		&messageStatus,
		&messagePaidAt,
	)
	if err != nil {
		t.Fatal(err)
	}

	if invoiceStatus != "paid" || messageStatus != "paid" {
		t.Fatalf("expected both records paid, got invoice=%s message=%s", invoiceStatus, messageStatus)
	}

	if storedPreimage != preimage {
		t.Fatalf("expected preimage %s, got %s", preimage, storedPreimage)
	}

	if invoicePaidAt.IsZero() || messagePaidAt.IsZero() {
		t.Fatal("expected paid timestamps to be set")
	}
}

func TestMarkPaidMissingInvoice(t *testing.T) {
	_, repository := newTestInvoiceRepository(t)

	if err := repository.MarkPaid(context.Background(), "missing", nil); err == nil {
		t.Fatal("expected error for missing invoice")
	}
}
