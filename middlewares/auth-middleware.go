package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/IdrisAkintobi/go-basic-crud/services"
	"github.com/IdrisAkintobi/go-basic-crud/utils"
	"github.com/jackc/pgx/v5"
)

type AuthData struct {
	SessionId int
	UserID    string
}

// NewAuthMiddleware returns a middleware function directly
func NewAuthMiddleware(db *pgx.Conn) func(http.Handler) http.Handler {
	ss := services.NewSessionService(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized - No token provided", http.StatusUnauthorized)
				return
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate token
			session, err := ss.FindSession(token)
			if err != nil && err != pgx.ErrNoRows {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if session == nil || session.ExpiresAt.Before(time.Now()) {
				http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
				return
			}

			// Store auth data in context
			authData := &AuthData{
				SessionId: session.ID,
				UserID:    session.UserId,
			}
			ctx := context.WithValue(r.Context(), utils.AuthUserCtxKey, authData)

			// Proceed with request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
