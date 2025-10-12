package auth

import (
	"errors"
	"net/http"
)

// Validator validates bearer tokens.
type Validator struct {
	secret string
}

// NewValidator creates a validator using the provided shared secret.
func NewValidator(secret string) *Validator {
	return &Validator{secret: secret}
}

// Middleware verifies Authorization headers on inbound HTTP requests.
func (v *Validator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1.- Extract the bearer token from the Authorization header.
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "missing authorization", http.StatusUnauthorized)
			return
		}

		// 2.- Compare against the configured secret before invoking the next handler.
		if !v.valid(token) {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ValidateToken returns an error when the provided token does not match the secret.
func (v *Validator) ValidateToken(token string) error {
	if !v.valid(token) {
		return errors.New("invalid token")
	}
	return nil
}

func (v *Validator) valid(token string) bool {
	return token == "Bearer "+v.secret
}
