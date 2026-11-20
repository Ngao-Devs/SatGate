package webhooks

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestWebhookHandler(t *testing.T) (*gin.Engine, *Service) {
	t.Helper()

	_, service := newTestWebhookService(t)
	handler := NewHandler(service)

	router := gin.New()
	router.POST("/webhooks/lightning", handler.Handle)

	return router, service
}

func TestHandlerProcessesPaidWebhook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router, _ := newTestWebhookHandler(t)
	body := webhookBody(t, "hash-123")

	request := httptest.NewRequest(
		http.MethodPost,
		"/webhooks/lightning",
		strings.NewReader(string(body)),
	)
	for key, values := range signedHeaders(body, "event-1") {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", response.Code, response.Body.String())
	}

	if !strings.Contains(response.Body.String(), `"status":"paid"`) {
		t.Fatalf("expected paid response, got %s", response.Body.String())
	}
}

func TestHandlerRejectsInvalidSignature(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router, _ := newTestWebhookHandler(t)
	body := webhookBody(t, "hash-123")

	request := httptest.NewRequest(
		http.MethodPost,
		"/webhooks/lightning",
		strings.NewReader(string(body)),
	)
	request.Header.Set("x-satgate-signature", "invalid")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}

	if !strings.Contains(response.Body.String(), "signature verification failed") {
		t.Fatalf("unexpected response: %s", response.Body.String())
	}
}

func TestHandlerRejectsInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router, _ := newTestWebhookHandler(t)
	body := []byte("{invalid-json")

	request := httptest.NewRequest(
		http.MethodPost,
		"/webhooks/lightning",
		strings.NewReader(string(body)),
	)
	for key, values := range signedHeaders(body, "event-1") {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}

	if !strings.Contains(response.Body.String(), "invalid webhook body") {
		t.Fatalf("unexpected response: %s", response.Body.String())
	}
}

func TestHandlerRejectsInvalidStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router, _ := newTestWebhookHandler(t)
	body := []byte(`{"invoice_id":"invoice-1","status":"pending"}`)

	request := httptest.NewRequest(
		http.MethodPost,
		"/webhooks/lightning",
		strings.NewReader(string(body)),
	)
	for key, values := range signedHeaders(body, "event-1") {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}

	if !strings.Contains(response.Body.String(), "webhook status is not paid") {
		t.Fatalf("unexpected response: %s", response.Body.String())
	}
}

func TestHandlerRejectsPaymentHashMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router, _ := newTestWebhookHandler(t)
	body := webhookBody(t, "wrong-hash")

	request := httptest.NewRequest(
		http.MethodPost,
		"/webhooks/lightning",
		strings.NewReader(string(body)),
	)
	for key, values := range signedHeaders(body, "event-1") {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}

	if !strings.Contains(response.Body.String(), "payment hash does not match invoice") {
		t.Fatalf("unexpected response: %s", response.Body.String())
	}
}
