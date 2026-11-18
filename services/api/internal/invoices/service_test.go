package invoices

import (
	"context"
	"testing"
)

func TestServiceGetStatus(t *testing.T) {
	_, repository := newTestInvoiceRepository(t)
	service := NewService(repository)

	response, err := service.GetStatus(context.Background(), "invoice-1")
	if err != nil {
		t.Fatal(err)
	}

	if response.InvoiceID != "invoice-1" {
		t.Fatalf("expected invoice-1, got %s", response.InvoiceID)
	}
}

func TestServiceMarkPaid(t *testing.T) {
	db, repository := newTestInvoiceRepository(t)
	service := NewService(repository)

	preimage := "service-test-preimage"
	if err := service.MarkPaid(context.Background(), "invoice-1", &preimage); err != nil {
		t.Fatal(err)
	}

	var status string
	if err := db.QueryRow(
		"SELECT status FROM invoices WHERE id = ?",
		"invoice-1",
	).Scan(&status); err != nil {
		t.Fatal(err)
	}

	if status != "paid" {
		t.Fatalf("expected paid, got %s", status)
	}
}
