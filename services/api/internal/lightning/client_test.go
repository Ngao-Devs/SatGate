package lightning

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateInvoiceSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.Header.Get("X-API-KEY") != "test-api-key" {
			t.Error("expected API key header")
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("expected JSON content type")
		}

		var payload struct {
			Variables struct {
				Input struct {
					WalletID string `json:"walletId"`
					Amount   int64  `json:"amount"`
					Memo     string `json:"memo"`
				} `json:"input"`
			} `json:"variables"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode request: %v", err)
		}

		if payload.Variables.Input.WalletID != "wallet-123" {
			t.Errorf("unexpected wallet ID: %s", payload.Variables.Input.WalletID)
		}
		if payload.Variables.Input.Amount != 250 {
			t.Errorf("unexpected amount: %d", payload.Variables.Input.Amount)
		}
		if payload.Variables.Input.Memo != "test payment" {
			t.Errorf("unexpected memo: %s", payload.Variables.Input.Memo)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": {
				"lnInvoiceCreate": {
					"invoice": {
						"paymentRequest": "lnbc-test",
						"paymentHash": "hash-123"
					},
					"errors": []
				}
			}
		}`))
	}))
	defer server.Close()

	result, err := createInvoice(
		server.Client(),
		server.URL,
		"test-api-key",
		"wallet-123",
		250,
		"test payment",
	)
	if err != nil {
		t.Fatalf("CreateInvoice returned error: %v", err)
	}

	if result.PaymentRequest != "lnbc-test" {
		t.Errorf("unexpected payment request: %s", result.PaymentRequest)
	}
	if result.PaymentHash != "hash-123" {
		t.Errorf("unexpected payment hash: %s", result.PaymentHash)
	}
}

func TestCreateInvoiceHandlesBlinkErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{
			"data": {
				"lnInvoiceCreate": {
					"invoice": {},
					"errors": [{"message": "insufficient balance"}]
				}
			}
		}`))
	}))
	defer server.Close()

	_, err := createInvoice(server.Client(), server.URL, "key", "wallet", 100, "memo")
	if err == nil || !strings.Contains(err.Error(), "insufficient balance") {
		t.Fatalf("expected Blink error, got %v", err)
	}
}

func TestCreateInvoiceRejectsEmptyInvoice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{
			"data": {
				"lnInvoiceCreate": {
					"invoice": {},
					"errors": []
				}
			}
		}`))
	}))
	defer server.Close()

	_, err := createInvoice(server.Client(), server.URL, "key", "wallet", 100, "memo")
	if err == nil || err.Error() != "blink returned empty invoice" {
		t.Fatalf("expected empty invoice error, got %v", err)
	}
}

func TestCreateInvoiceRejectsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	_, err := createInvoice(server.Client(), server.URL, "key", "wallet", 100, "memo")
	if err == nil || !strings.Contains(err.Error(), "503") {
		t.Fatalf("expected HTTP 503 error, got %v", err)
	}
}

func TestCreateInvoiceRejectsInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not-json"))
	}))
	defer server.Close()

	_, err := createInvoice(server.Client(), server.URL, "key", "wallet", 100, "memo")
	if err == nil {
		t.Fatal("expected JSON decoding error")
	}
}
