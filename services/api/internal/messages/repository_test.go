package messages

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func newTestRepository(t *testing.T) *Repository {
	t.Helper()

	db, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "messages.db"))
	if err != nil {
		t.Fatal(err)
	}

	db.SetMaxOpenConns(1)

	schema := `
		CREATE TABLE forms (
			id TEXT PRIMARY KEY,
			owner_id TEXT NOT NULL
		);

		CREATE TABLE messages (
			id TEXT PRIMARY KEY,
                        owner_id TEXT NOT NULL,
			form_id TEXT NOT NULL,
			sender_name TEXT NOT NULL,
			sender_email TEXT NOT NULL,
			body TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			paid_at DATETIME
		);

		CREATE TABLE invoices (
			id TEXT PRIMARY KEY,
			message_id TEXT NOT NULL,
			amount_sats INTEGER NOT NULL
		);
	`

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return NewRepository(db)
}

func TestListPaidByOwner(t *testing.T) {
	repository := newTestRepository(t)
	ctx := context.Background()
	now := time.Now().UTC()

	_, err := repository.db.Exec(`
		INSERT INTO forms (id, owner_id) VALUES
			('form-a', 'owner-a'),
			('form-b', 'owner-b');

		INSERT INTO messages (
			id, owner_id, form_id, sender_name, sender_email,
			body, status, created_at, paid_at
		) VALUES
			('message-a1', 'owner-a', 'form-a', 'Alice', 'alice@example.com',
			 'First message', 'paid', ?, ?),
			('message-a2', 'owner-a', 'form-a', 'Bob', 'bob@example.com',
			 'Pending message', 'pending', ?, NULL),
			('message-b1', 'owner-b', 'form-b', 'Charlie', 'charlie@example.com',
			 'Other owner message', 'paid', ?, ?);

		INSERT INTO invoices (id, message_id, amount_sats) VALUES
			('invoice-a1', 'message-a1', 100),
			('invoice-b1', 'message-b1', 200);
	`,
		now, now,
		now, now,
		now, now,
	)
	if err != nil {
		t.Fatal(err)
	}

	messages, err := repository.ListPaidByOwner(ctx, "owner-a")
	if err != nil {
		t.Fatal(err)
	}

	if len(messages) != 1 {
		t.Fatalf("expected 1 paid message, got %d", len(messages))
	}

	message := messages[0]

	if message.ID != "message-a1" {
		t.Errorf("expected message-a1, got %q", message.ID)
	}

	if message.FormID != "form-a" {
		t.Errorf("expected form-a, got %q", message.FormID)
	}

	if message.AmountSats != 100 {
		t.Errorf("expected amount 100 sats, got %d", message.AmountSats)
	}

	if message.Status != "paid" {
		t.Errorf("expected paid status, got %q", message.Status)
	}

	if message.PaidAt == nil {
		t.Error("expected paid_at to be populated")
	}
}

func TestListPaidByOwnerReturnsEmptySlice(t *testing.T) {
	repository := newTestRepository(t)

	messages, err := repository.ListPaidByOwner(context.Background(), "unknown-owner")
	if err != nil {
		t.Fatal(err)
	}

	if messages == nil {
		t.Fatal("expected an empty slice, got nil")
	}

	if len(messages) != 0 {
		t.Fatalf("expected 0 messages, got %d", len(messages))
	}
}
