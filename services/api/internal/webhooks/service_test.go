package webhooks

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"testing"

	"satgate/internal/invoices"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const testWebhookSecret = "test-secret"

func newTestWebhookService(t *testing.T) (*sql.DB, *Service) {
	t.Helper()

	db, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "webhooks.db"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE forms (
			id TEXT PRIMARY KEY,
			sat_balance INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE messages (
			id TEXT PRIMARY KEY,
			form_id TEXT NOT NULL,
			status TEXT NOT NULL,
			paid_at DATETIME
		);

		CREATE TABLE invoices (
			id TEXT PRIMARY KEY,
			message_id TEXT NOT NULL,
			amount_sats INTEGER NOT NULL,
			payment_hash TEXT,
			status TEXT NOT NULL,
			preimage TEXT,
			paid_at DATETIME
		);

		CREATE TABLE webhook_events (
			id TEXT PRIMARY KEY,
			invoice_id TEXT NOT NULL,
			payload TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);

		INSERT INTO forms (id, sat_balance) VALUES ('form-1', 0);
		INSERT INTO messages (id, form_id, status)
		VALUES ('message-1', 'form-1', 'pending');
		INSERT INTO invoices (id, message_id, amount_sats, payment_hash, status)
		VALUES ('invoice-1', 'message-1', 25, 'hash-123', 'pending');
	`)
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	invoiceRepository := invoices.NewRepository(db)
	invoiceService := invoices.NewService(invoiceRepository)
	webhookRepository := NewRepository(db)

	return db, NewService(webhookRepository, invoiceService, testWebhookSecret)
}

func signedHeaders(body []byte, eventID string) http.Header {
	mac := hmac.New(sha256.New, []byte(testWebhookSecret))
	mac.Write(body)

	headers := http.Header{}
	headers.Set("x-satgate-signature", hex.EncodeToString(mac.Sum(nil)))
	headers.Set("x-satgate-event-id", eventID)

	return headers
}

func webhookBody(t *testing.T, paymentHash string) []byte {
	t.Helper()

	payload := LightningWebhook{
		InvoiceID:   "invoice-1",
		PaymentHash: &paymentHash,
		Preimage:    nil,
		Status:      "paid",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	return body
}

func TestProcessMarksInvoicePaid(t *testing.T) {
	db, service := newTestWebhookService(t)

	body := webhookBody(t, "hash-123")
	response, err := service.Process(
		context.Background(),
		signedHeaders(body, "event-1"),
		body,
	)
	if err != nil {
		t.Fatal(err)
	}

	if response.InvoiceID != "invoice-1" || response.Status != "paid" {
		t.Fatalf("unexpected response: %+v", response)
	}

	var messageStatus string
	if err := db.QueryRow(
		"SELECT status FROM messages WHERE id = ?",
		"message-1",
	).Scan(&messageStatus); err != nil {
		t.Fatal(err)
	}

	if messageStatus != "paid" {
		t.Fatalf("expected message paid, got %s", messageStatus)
	}
}

func TestProcessDuplicateEvent(t *testing.T) {
	_, service := newTestWebhookService(t)

	body := webhookBody(t, "hash-123")
	headers := signedHeaders(body, "event-1")

	if _, err := service.Process(context.Background(), headers, body); err != nil {
		t.Fatal(err)
	}

	response, err := service.Process(context.Background(), headers, body)
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != "paid" {
		t.Fatalf("expected paid status, got %s", response.Status)
	}
}

func TestProcessRejectsInvalidSignature(t *testing.T) {
	_, service := newTestWebhookService(t)

	body := webhookBody(t, "hash-123")
	headers := http.Header{}
	headers.Set("x-satgate-signature", "invalid")

	_, err := service.Process(context.Background(), headers, body)

	var processingError *ProcessingError
	if err == nil || !errors.As(err, &processingError) || processingError.Kind != ErrSignature {
		t.Fatalf("expected signature error, got %v", err)
	}
}

func TestProcessRejectsInvalidStatus(t *testing.T) {
	_, service := newTestWebhookService(t)

	payload := LightningWebhook{
		InvoiceID: "invoice-1",
		Status:    "pending",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.Process(
		context.Background(),
		signedHeaders(body, "event-1"),
		body,
	)

	var processingError *ProcessingError
	if err == nil || !errors.As(err, &processingError) || processingError.Kind != ErrInvalidStatus {
		t.Fatalf("expected invalid status error, got %v", err)
	}
}

func TestProcessRejectsPaymentHashMismatch(t *testing.T) {
	_, service := newTestWebhookService(t)

	body := webhookBody(t, "wrong-hash")
	_, err := service.Process(
		context.Background(),
		signedHeaders(body, "event-1"),
		body,
	)

	var processingError *ProcessingError
	if err == nil || !errors.As(err, &processingError) || processingError.Kind != ErrPaymentHash {
		t.Fatalf("expected payment hash error, got %v", err)
	}
}
