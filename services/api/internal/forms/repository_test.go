package forms

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func newFormsTestRepository(t *testing.T) *Repository {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })

	_, err = db.Exec(`
		CREATE TABLE forms (
			id TEXT PRIMARY KEY,
			owner_id TEXT NOT NULL,
			name TEXT NOT NULL,
			domain TEXT NOT NULL,
			amount_sats INTEGER NOT NULL,
			created_at DATETIME NOT NULL
		);
		CREATE TABLE messages (id TEXT PRIMARY KEY, form_id TEXT NOT NULL);
		CREATE TABLE invoices (id TEXT PRIMARY KEY, message_id TEXT NOT NULL);
		CREATE TABLE webhook_events (id TEXT PRIMARY KEY, invoice_id TEXT NOT NULL);
	`)
	if err != nil {
		t.Fatal(err)
	}

	return NewRepository(db)
}

func TestRepositoryListByOwner(t *testing.T) {
	repo := newFormsTestRepository(t)
	ctx := context.Background()
	now := time.Now().UTC()

	for _, fixture := range []struct {
		id, owner, name string
	}{
		{"form-a", "owner-a", "A form"},
		{"form-b", "owner-b", "B form"},
	} {
		err := repo.Create(ctx, fixture.owner, Form{
			ID: fixture.id, Name: fixture.name, Domain: "example.com",
			AmountSats: 100, CreatedAt: now,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	got, err := repo.ListByOwner(ctx, "owner-a")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "form-a" {
		t.Fatalf("got %+v; want only form-a", got)
	}
}

func TestRepositoryGetByIDAndOwner(t *testing.T) {
	repo := newFormsTestRepository(t)
	ctx := context.Background()

	err := repo.Create(ctx, "owner-a", Form{
		ID: "form-a", Name: "A form", Domain: "example.com",
		AmountSats: 100, CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := repo.GetByIDAndOwner(ctx, "form-a", "owner-b"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got error %v; want ErrNotFound", err)
	}

	got, err := repo.GetByIDAndOwner(ctx, "form-a", "owner-a")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "form-a" {
		t.Fatalf("got ID %q; want form-a", got.ID)
	}
}

func TestRepositoryUpdateIsOwnerScoped(t *testing.T) {
	repo := newFormsTestRepository(t)
	ctx := context.Background()

	err := repo.Create(ctx, "owner-a", Form{
		ID: "form-a", Name: "Original", Domain: "example.com",
		AmountSats: 100, CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		t.Fatal(err)
	}

	updated := Form{
		ID: "form-a", Name: "Updated",
		Domain: "new.example.com", AmountSats: 200,
	}

	if err := repo.Update(ctx, "owner-b", updated); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got error %v; want ErrNotFound", err)
	}

	got, err := repo.GetByIDAndOwner(ctx, "form-a", "owner-a")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Original" || got.AmountSats != 100 {
		t.Fatalf("unauthorized update changed form: %+v", got)
	}

	if err := repo.Update(ctx, "owner-a", updated); err != nil {
		t.Fatal(err)
	}
	got, err = repo.GetByIDAndOwner(ctx, "form-a", "owner-a")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Updated" || got.AmountSats != 200 {
		t.Fatalf("got %+v; want updated values", got)
	}
}

func insertRelatedFixtures(t *testing.T, repo *Repository) {
	t.Helper()

	_, err := repo.db.Exec(`
		INSERT INTO forms VALUES ('form-a', 'owner-a', 'A form', 'example.com', 100, ?);
		INSERT INTO messages VALUES ('message-a', 'form-a');
		INSERT INTO invoices VALUES ('invoice-a', 'message-a');
		INSERT INTO webhook_events VALUES ('event-a', 'invoice-a');
	`, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
}

func assertRelatedFixturesExist(t *testing.T, repo *Repository, want int) {
	t.Helper()

	for table, id := range map[string]string{
		"forms":          "form-a",
		"messages":       "message-a",
		"invoices":       "invoice-a",
		"webhook_events": "event-a",
	} {
		var count int
		if err := repo.db.QueryRow("SELECT COUNT(*) FROM "+table+" WHERE id = ?", id).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != want {
			t.Errorf("%s count = %d; want %d", table, count, want)
		}
	}
}

func TestRepositoryDeleteRemovesRelatedRecords(t *testing.T) {
	repo := newFormsTestRepository(t)
	insertRelatedFixtures(t, repo)

	if err := repo.Delete(context.Background(), "form-a", "owner-a"); err != nil {
		t.Fatal(err)
	}
	assertRelatedFixturesExist(t, repo, 0)
}

func TestRepositoryDeleteRollsBackOnFailure(t *testing.T) {
	repo := newFormsTestRepository(t)
	insertRelatedFixtures(t, repo)

	_, err := repo.db.Exec(`
		CREATE TRIGGER reject_invoice_delete
		BEFORE DELETE ON invoices
		BEGIN
			SELECT RAISE(ABORT, 'delete blocked for test');
		END;
	`)
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.Delete(context.Background(), "form-a", "owner-a"); err == nil {
		t.Fatal("expected Delete to fail")
	}

	assertRelatedFixturesExist(t, repo, 1)
}
