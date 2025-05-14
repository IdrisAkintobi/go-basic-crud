package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/IdrisAkintobi/go-basic-crud/database"
	"github.com/IdrisAkintobi/go-basic-crud/handlers"
	"github.com/IdrisAkintobi/go-basic-crud/middlewares"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func gracefulShutdown(db *pgxpool.Pool) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	// Wait for signal
	<-ch
	signal.Stop(ch)

	// Create context that timeout in expected duration for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Do all necessary cleanup
	database.DisconnectDB(ctx, db)

	// Exit process
	os.Exit(0)
}

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3003"
	}

	// Connect database
	conn, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	// Activate graceful shutdown
	go gracefulShutdown(conn)

	// Setup handler
	uh := handlers.NewUserHandler(conn)
	ah := handlers.NewAuthHandler(conn)

	//create auth middleware
	authMiddleware := middlewares.NewAuthMiddleware(conn)

	// Setup routers
	r := chi.NewRouter()
	r.Use(httprate.LimitByIP(100, 1*time.Minute), middleware.CleanPath, middleware.StripSlashes, middleware.Logger, middleware.Recoverer)
	r.Post("/register", uh.RegisterUser)
	r.With(middlewares.GetUserFingerprint).Post("/login", ah.Login)
	r.With(authMiddleware.Register()).Get("/whoami", ah.WhoAmI)
	r.With(authMiddleware.Register()).Get("/active-sessions", ah.GetActiveSessions)
	r.With(authMiddleware.Register()).Post("/logout", ah.LogOut)

	// Start server
	fmt.Printf("Server starting on %v\n", PORT)
	http.ListenAndServe(":3003", r)
}
