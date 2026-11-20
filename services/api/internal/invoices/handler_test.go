package invoices

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestInvoiceHandler(t *testing.T) (*gin.Engine, *Repository) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	_, repository := newTestInvoiceRepository(t)
	service := NewService(repository)
	handler := NewHandler(service)

	router := gin.New()
	router.GET("/invoices/:id", handler.Get)
	router.POST("/invoices/:id/mock-pay", handler.MockPay)

	return router, repository
}

func TestHandlerGetInvoice(t *testing.T) {
	router, _ := newTestInvoiceHandler(t)

	request := httptest.NewRequest(http.MethodGet, "/invoices/invoice-1", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
}

func TestHandlerGetMissingInvoice(t *testing.T) {
	router, _ := newTestInvoiceHandler(t)

	request := httptest.NewRequest(http.MethodGet, "/invoices/missing", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", response.Code)
	}
}

func TestHandlerMockPay(t *testing.T) {
	router, repository := newTestInvoiceHandler(t)

	request := httptest.NewRequest(http.MethodPost, "/invoices/invoice-1/mock-pay", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", response.Code, response.Body.String())
	}

	status, err := repository.GetStatus(request.Context(), "invoice-1")
	if err != nil {
		t.Fatal(err)
	}

	if status.Status != "paid" {
		t.Fatalf("expected paid, got %s", status.Status)
	}
}
