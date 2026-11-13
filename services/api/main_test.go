package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"satgate/internal/middleware"
)

func newTestAPI(t *testing.T) (*sql.DB, *gin.Engine) {
	t.Helper()

	db, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}

	if err := runMigrations(db); err != nil {
		db.Close()
		t.Fatal(err)
	}

	now := time.Now().UTC()

	forms := []struct {
		id      string
		ownerID string
		name    string
	}{
		{"form-owner-a", "owner-a", "Owner A Form"},
		{"form-owner-b", "owner-b", "Owner B Form"},
	}

	for _, form := range forms {
		_, err := db.Exec(`
			INSERT INTO forms (id, owner_id, name, domain, amount_sats, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, form.id, form.ownerID, form.name, "localhost", 5, now)

		if err != nil {
			db.Close()
			t.Fatal(err)
		}
	}

	for _, formID := range []string{"form-owner-a", "form-owner-b"} {
		_, err := db.Exec(`
			INSERT INTO messages (
				id, form_id, sender_name, sender_email,
				body, status, created_at, paid_at
			)
			VALUES (?, ?, ?, ?, ?, 'paid', ?, ?)
		`,
			"message-"+formID,
			formID,
			"Test Sender",
			"sender@example.com",
			"Test message",
			now,
			now,
		)

		if err != nil {
			db.Close()
			t.Fatal(err)
		}
	}

	state := &AppState{DB: db}

	router := gin.New()
	protected := router.Group("/api")
	protected.Use(middleware.Auth("", true))

	protected.GET("/forms", state.listForms)
	protected.PUT("/forms/:id", state.updateForm)
	protected.DELETE("/forms/:id", state.deleteForm)
	protected.GET("/messages", state.listMessages)

	t.Cleanup(func() {
		db.Close()
	})

	return db, router
}

func performRequest(
	router *gin.Engine,
	method string,
	path string,
	body string,
	token string,
) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder
}

func TestOwnerCannotAccessAnotherOwnersForm(t *testing.T) {
	db, router := newTestAPI(t)

	// Owner A must not see Owner B's form.
	response := performRequest(
		router, http.MethodGet, "/api/forms", "", "mock-owner-owner-a",
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}

	var forms []Form
	if err := json.Unmarshal(response.Body.Bytes(), &forms); err != nil {
		t.Fatal(err)
	}

	if len(forms) != 1 || forms[0].ID != "form-owner-a" {
		t.Fatalf("owner A received unexpected forms: %+v", forms)
	}

	// Owner A must not update Owner B's form.
	response = performRequest(
		router,
		http.MethodPut,
		"/api/forms/form-owner-b",
		`{"name":"Hijacked","domain":"localhost","amount_sats":10}`,
		"mock-owner-owner-a",
	)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected 404 updating another owner's form, got %d", response.Code)
	}

	// Owner A must not delete Owner B's form.
	response = performRequest(
		router,
		http.MethodDelete,
		"/api/forms/form-owner-b",
		"",
		"mock-owner-owner-a",
	)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected 404 deleting another owner's form, got %d", response.Code)
	}

	// Confirm Owner B's form still exists and was not modified.
	var name string
	err := db.QueryRow(
		"SELECT name FROM forms WHERE id = ?",
		"form-owner-b",
	).Scan(&name)

	if err != nil {
		t.Fatalf("owner B's form should still exist: %v", err)
	}

	if name != "Owner B Form" {
		t.Fatalf("owner B's form was modified: %s", name)
	}
}

func TestOwnerCanUpdateOwnForm(t *testing.T) {
	_, router := newTestAPI(t)

	response := performRequest(
		router,
		http.MethodPut,
		"/api/forms/form-owner-a",
		`{"name":"Updated Form","domain":"localhost","amount_sats":10}`,
		"mock-owner-owner-a",
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200 updating own form, got %d: %s",
			response.Code, response.Body.String())
	}

	var form Form
	if err := json.Unmarshal(response.Body.Bytes(), &form); err != nil {
		t.Fatal(err)
	}

	if form.Name != "Updated Form" {
		t.Fatalf("expected updated form name, got %q", form.Name)
	}
}

func TestOwnerOnlySeesOwnPaidMessages(t *testing.T) {
	_, router := newTestAPI(t)

	response := performRequest(
		router,
		http.MethodGet,
		"/api/messages",
		"",
		"mock-owner-owner-a",
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}

	var messages []Message
	if err := json.Unmarshal(response.Body.Bytes(), &messages); err != nil {
		t.Fatal(err)
	}

	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}

	if messages[0].FormID != "form-owner-a" {
		t.Fatalf("owner A received another owner's message: %+v", messages[0])
	}
}
