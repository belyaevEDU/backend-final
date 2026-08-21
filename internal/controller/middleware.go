package controller

import (
	"context"
	"net/http"
	"strings"

	"github.com/belyaevedu/remote-code-service/internal/controller/handlers"
	"github.com/belyaevedu/remote-code-service/internal/port"
)

const (
	bearerScheme = "Bearer "
	// context key for the authenticated user id, stored for downstream handlers
	UserIDKey = "user_id"
)

func AuthMiddleware(auth port.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				unauthorizedResponseHelper(w, "missing Authorization header")
				return
			}

			if !strings.HasPrefix(header, bearerScheme) {
				unauthorizedResponseHelper(w, "unsupported Authorization scheme")
				return
			}

			token := strings.TrimSpace(header[len(bearerScheme):])
			if token == "" {
				unauthorizedResponseHelper(w, "empty bearer token")
				return
			}

			userID, err := auth.Authenticate(token)
			if err != nil {
				unauthorizedResponseHelper(w, "invalid or expired session")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func unauthorizedResponseHelper(w http.ResponseWriter, msg string) {
	handlers.WriteJSON(w, http.StatusUnauthorized, handlers.ErrorResponse{Error: msg})
}
