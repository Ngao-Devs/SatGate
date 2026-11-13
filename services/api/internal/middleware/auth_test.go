package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret-for-satgate"

func setupRouter(secret string, allowMock bool) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET(
		"/protected",
		Auth(secret, allowMock),
		func(c *gin.Context) {
			ownerID, ok := OwnerID(c)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "missing owner"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"owner_id": ownerID})
		},
	)

	return router
}

func createToken(t *testing.T, secret string, claims Claims) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	return signed
}

func request(router *gin.Engine, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder
}

func TestAuthAcceptsValidJWT(t *testing.T) {
	router := setupRouter(testSecret, false)

	token := createToken(t, testSecret, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "owner-123",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})

	response := request(router, token)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	if response.Body.String() != `{"owner_id":"owner-123"}` {
		t.Fatalf("unexpected response body: %s", response.Body.String())
	}
}

func TestAuthRejectsMissingToken(t *testing.T) {
	router := setupRouter(testSecret, false)

	response := request(router, "")

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}

func TestAuthRejectsInvalidToken(t *testing.T) {
	router := setupRouter(testSecret, false)

	response := request(router, "invalid-token")

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}

func TestAuthRejectsExpiredToken(t *testing.T) {
	router := setupRouter(testSecret, false)

	token := createToken(t, testSecret, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "owner-123",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	})

	response := request(router, token)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}

func TestAuthRejectsWrongSecret(t *testing.T) {
	router := setupRouter(testSecret, false)

	token := createToken(t, "different-secret", Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "owner-123",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})

	response := request(router, token)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}

func TestAuthAcceptsMockTokenInDevelopment(t *testing.T) {
	router := setupRouter("", true)

	response := request(router, "mock-owner-demo-user")

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	if response.Body.String() != `{"owner_id":"demo-user"}` {
		t.Fatalf("unexpected response body: %s", response.Body.String())
	}
}

func TestAuthRejectsMockTokenInProduction(t *testing.T) {
	router := setupRouter(testSecret, false)

	response := request(router, "mock-owner-demo-user")

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}

func TestAuthRejectsEmptyMockOwner(t *testing.T) {
	router := setupRouter("", true)

	response := request(router, "mock-owner-")

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}

func TestAuthRejectsMissingJWTSubject(t *testing.T) {
	router := setupRouter(testSecret, false)

	token := createToken(t, testSecret, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})

	response := request(router, token)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}
