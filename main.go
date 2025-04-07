package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"

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
	r.With(authMiddleware.Register()).Post("/logout", ah.LogOut)

	// Start server
	fmt.Printf("Server starting on %v\n", PORT)
	http.ListenAndServe(":3003", r)
}
