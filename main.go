package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/IdrisAkintobi/go-basic-crud/database"
	"github.com/IdrisAkintobi/go-basic-crud/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

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

	// Setup routers
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", uh.RegisterUser)

	// Start server
	fmt.Printf("Server starting on %v\n", PORT)
	http.ListenAndServe(":3003", r)
}
