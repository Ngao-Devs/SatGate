package messages

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"satgate/internal/middleware"
)

func TestHandlerListRequiresAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repository := newTestRepository(t)
	service := NewService(repository)
	handler := NewHandler(service)

	router := gin.New()
	protected := router.Group("/api")
	protected.Use(middleware.Auth("", true))
	protected.GET("/messages", handler.List)

	request := httptest.NewRequest(http.MethodGet, "/api/messages", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}

func TestHandlerListReturnsOnlyOwnersPaidMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repository := newTestRepository(t)
	now := time.Now().UTC()

	_, err := repository.db.Exec(`
		INSERT INTO forms (id, owner_id) VALUES
			('form-a', 'owner-a'),
			('form-b', 'owner-b');

		INSERT INTO messages (
			id, owner_id, form_id, sender_name, sender_email,
			body, status, created_at, paid_at
		) VALUES
			('message-a', 'owner-a', 'form-a', 'Alice', 'alice@example.com',
			 'Owner A message', 'paid', ?, ?),
			('message-b', 'owner-b', 'form-b', 'Bob', 'bob@example.com',
			 'Owner B message', 'paid', ?, ?);

		INSERT INTO invoices (id, message_id, amount_sats) VALUES
			('invoice-a', 'message-a', 100),
			('invoice-b', 'message-b', 200);
	`,
		now, now,
		now, now,
	)
	if err != nil {
		t.Fatal(err)
	}

	service := NewService(repository)
	handler := NewHandler(service)

	router := gin.New()
	protected := router.Group("/api")
	protected.Use(middleware.Auth("", true))
	protected.GET("/messages", handler.List)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/messages",
		strings.NewReader(""),
	)
	request.Header.Set("Authorization", "Bearer mock-owner-owner-a")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", response.Code, response.Body.String())
	}

	if !strings.Contains(response.Body.String(), "message-a") {
		t.Fatal("expected owner A's message in response")
	}

	if strings.Contains(response.Body.String(), "message-b") {
		t.Fatal("owner A must not receive owner B's message")
	}
}
