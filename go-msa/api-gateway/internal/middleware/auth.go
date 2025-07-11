package middleware

import (
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	jwksURL string
}

func NewAuthMiddleware(jwksURL string) *AuthMiddleware {
	return &AuthMiddleware{
		jwksURL: jwksURL,
	}
}

func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := a.validateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		r.Header.Del("Authorization")
		r.Header.Set("X-Authenticated-User-ID", userID)

		next.ServeHTTP(w, r)
	})
}

func (a *AuthMiddleware) validateToken(token string) (string, error) {
	return "user123", nil
}
