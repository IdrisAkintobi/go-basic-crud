package middlewares

import (
	"context"
	"net"
	"net/http"

	"github.com/IdrisAkintobi/go-basic-crud/utils"
)

type UserFingerprint struct {
	IPAddress string
	UserAgent string
	DeviceId  string
}

// GetUserFingerprint is a middleware that captures user identification details
func GetUserFingerprint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Split host and port
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr // fallback just in case
		}

		fingerprint := &UserFingerprint{
			IPAddress: ip,
			UserAgent: r.UserAgent(),
			DeviceId:  r.Header.Get("X-DeviceID"),
		}

		ctx := context.WithValue(r.Context(), utils.FPCtxKey, fingerprint)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
