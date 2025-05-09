package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/IdrisAkintobi/go-basic-crud/services"
	"github.com/IdrisAkintobi/go-basic-crud/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthData struct {
	SessionId int
	UserID    string
}

type AuthMiddleware struct {
	ss *services.SessionService
}

func NewAuthMiddleware(db *pgxpool.Pool) *AuthMiddleware {
	return &AuthMiddleware{
		ss: services.NewSessionService(db),
	}
}

// NewAuthMiddleware returns a middleware function directly
func (um *AuthMiddleware) Register() func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.SendErrorResponse(w, "no token provided", http.StatusUnauthorized)
				return
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate token
			session, err := um.ss.FindSession(token)
			if err != nil && err != pgx.ErrNoRows {
				utils.SendErrorResponse(w, "internal server error", http.StatusInternalServerError)
				return
			}
			if session == nil || session.ExpiresAt.Before(time.Now()) {
				utils.SendErrorResponse(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// Check if session needs to be refreshed
			// If expiration is within the refresh window, extend the session
			timeRemaining := time.Until(session.ExpiresAt)
			if timeRemaining < um.ss.SessionRefreshWindow {
				_ = um.ss.ExtendSession(token)
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
