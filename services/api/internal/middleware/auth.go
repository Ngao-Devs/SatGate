package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const OwnerIDKey = "owner_id"

type Claims struct {
	jwt.RegisteredClaims
}

func Auth(secret string, allowMockTokens bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.Fields(header)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			unauthorized(c)
			return
		}

		tokenString := parts[1]

		if strings.HasPrefix(tokenString, "mock-owner-") {
			if !allowMockTokens {
				unauthorized(c)
				return
			}

			ownerID := strings.TrimPrefix(tokenString, "mock-owner-")
			if ownerID == "" {
				unauthorized(c)
				return
			}

			c.Set(OwnerIDKey, ownerID)
			c.Next()
			return
		}

		if secret == "" {
			unauthorized(c)
			return
		}

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		)

		if err != nil || !token.Valid ||
			claims.Subject == "" ||
			claims.ExpiresAt == nil {
			unauthorized(c)
			return
		}

		c.Set(OwnerIDKey, claims.Subject)
		c.Next()
	}
}

func OwnerID(c *gin.Context) (string, bool) {
	value, exists := c.Get(OwnerIDKey)
	if !exists {
		return "", false
	}

	id, ok := value.(string)
	return id, ok && id != ""
}

func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "unauthorized",
	})
}
