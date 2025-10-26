package routes

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/IdrisAkintobi/go-basic-crud/handlers"
	"github.com/IdrisAkintobi/go-basic-crud/middlewares"
)

func SetupRoutes(conn *pgxpool.Pool) *chi.Mux {
	// Setup handlers
	uh := handlers.NewUserHandler(conn)
	ah := handlers.NewAuthHandler(conn)

	// Create auth middleware
	authMiddleware := middlewares.NewAuthMiddleware(conn)

	// Setup router
	r := chi.NewRouter()

	// Global middleware
	r.Use(
		httprate.LimitByIP(100, 1*time.Minute),
		middleware.CleanPath,
		middleware.StripSlashes,
		middleware.Logger,
		middleware.Recoverer,
		middlewares.GetUserFingerprint,
	)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Post("/register", uh.RegisterUser)
		r.Post("/login", ah.Login)

		// Protected routes group
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Register())
			r.Get("/whoami", ah.WhoAmI)
			r.Get("/active-sessions", ah.GetActiveSessions)
			r.Post("/logout", ah.LogOut)
		})
	})

	return r
}
