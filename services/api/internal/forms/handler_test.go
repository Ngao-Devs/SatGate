package forms

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"satgate/internal/middleware"

	"github.com/gin-gonic/gin"
)

func newTestHandler(t *testing.T) (*gin.Engine, *Service) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	service := newFormsTestService(t)
	handler := NewHandler(service)

	router := gin.New()
	router.Use(middleware.Auth("", true))

	router.GET("/forms", handler.List)
	router.POST("/forms", handler.Create)
	router.PUT("/forms/:id", handler.Update)
	router.DELETE("/forms/:id", handler.Delete)

	return router, service
}

func sendRequest(router http.Handler, method, path, token, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")

	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	return recorder
}

func createHandlerTestForm(t *testing.T, service *Service, ownerID string) *Form {
	t.Helper()

	form, err := service.Create(context.Background(), ownerID, CreateFormRequest{
		Name:       "Test form",
		Domain:     "example.com",
		AmountSats: 100,
	})
	if err != nil {
		t.Fatal(err)
	}

	return form
}

func TestHandlerListForms(t *testing.T) {
	router, service := newTestHandler(t)

	createHandlerTestForm(t, service, "owner-a")
	createHandlerTestForm(t, service, "owner-b")

	response := sendRequest(router, http.MethodGet, "/forms", "mock-owner-owner-a", "")

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", response.Code, http.StatusOK)
	}

	var forms []Form
	if err := json.Unmarshal(response.Body.Bytes(), &forms); err != nil {
		t.Fatal(err)
	}

	if len(forms) != 1 {
		t.Fatalf("got %d forms; want 1", len(forms))
	}
}

func TestHandlerCreateForm(t *testing.T) {
	router, _ := newTestHandler(t)

	response := sendRequest(
		router,
		http.MethodPost,
		"/forms",
		"mock-owner-owner-a",
		`{"name":"Contact form","domain":"example.com","amount_sats":100}`,
	)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d. Body: %s", response.Code, http.StatusOK, response.Body.String())
	}

	var form Form
	if err := json.Unmarshal(response.Body.Bytes(), &form); err != nil {
		t.Fatal(err)
	}

	if form.Name != "Contact form" || form.Domain != "example.com" {
		t.Fatalf("unexpected response: %+v", form)
	}
	if form.ID == "" {
		t.Fatal("expected generated form ID")
	}
}

func TestHandlerCreateFormValidation(t *testing.T) {
	router, _ := newTestHandler(t)

	response := sendRequest(
		router,
		http.MethodPost,
		"/forms",
		"mock-owner-owner-a",
		`{"name":"Contact form","domain":" ","amount_sats":100}`,
	)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d; want %d", response.Code, http.StatusBadRequest)
	}

	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}

	if body["error"] != "allowed domain is required" {
		t.Fatalf("error = %q; want allowed domain is required", body["error"])
	}
}

func TestHandlerRequiresAuthentication(t *testing.T) {
	router, _ := newTestHandler(t)

	response := sendRequest(router, http.MethodGet, "/forms", "", "")

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d; want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestHandlerCannotUpdateOrDeleteAnotherOwnersForm(t *testing.T) {
	router, service := newTestHandler(t)
	form := createHandlerTestForm(t, service, "owner-a")

	updateResponse := sendRequest(
		router,
		http.MethodPut,
		"/forms/"+form.ID,
		"mock-owner-owner-b",
		`{"name":"Hijacked","domain":"attacker.example.com","amount_sats":200}`,
	)

	if updateResponse.Code != http.StatusNotFound {
		t.Fatalf("update status = %d; want %d", updateResponse.Code, http.StatusNotFound)
	}

	deleteResponse := sendRequest(
		router,
		http.MethodDelete,
		"/forms/"+form.ID,
		"mock-owner-owner-b",
		"",
	)

	if deleteResponse.Code != http.StatusNotFound {
		t.Fatalf("delete status = %d; want %d", deleteResponse.Code, http.StatusNotFound)
	}
}

func TestHandlerDeleteForm(t *testing.T) {
	router, service := newTestHandler(t)
	form := createHandlerTestForm(t, service, "owner-a")

	response := sendRequest(
		router,
		http.MethodDelete,
		"/forms/"+form.ID,
		"mock-owner-owner-a",
		"",
	)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d; want %d", response.Code, http.StatusNoContent)
	}

	_, err := service.GetPublic(context.Background(), form.ID)
	if err == nil {
		t.Fatal("expected deleted form not to exist")
	}
}

func TestHandlerUpdateForm(t *testing.T) {
	router, service := newTestHandler(t)
	form := createHandlerTestForm(t, service, "owner-a")

	response := sendRequest(
		router,
		http.MethodPut,
		"/forms/"+form.ID,
		"mock-owner-owner-a",
		`{"name":"Updated form","domain":"updated.example.com","amount_sats":250}`,
	)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d. Body: %s", response.Code, http.StatusOK, response.Body.String())
	}

	var updated Form
	if err := json.Unmarshal(response.Body.Bytes(), &updated); err != nil {
		t.Fatal(err)
	}

	if updated.Name != "Updated form" || updated.AmountSats != 250 {
		t.Fatalf("unexpected updated form: %+v", updated)
	}

	_ = time.Now()
}
