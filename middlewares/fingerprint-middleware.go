package middlewares

import (
	"context"
	"net/http"

	"github.com/IdrisAkintobi/go-basic-crud/utils"
)

type UserFingerprint struct {
	IPAddress string
	UserAgent string
}

// GetUserFingerprint is a middleware that captures user identification details
func GetUserFingerprint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Create a new UserFingerprint instance
		fingerprint := &UserFingerprint{
			IPAddress: r.RemoteAddr,
			UserAgent: r.UserAgent(),
		}

		// Add the fingerprint to the request context
		ctx := context.WithValue(r.Context(), utils.FPCtxKey, fingerprint)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
