package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/IdrisAkintobi/go-basic-crud/config"
	"github.com/IdrisAkintobi/go-basic-crud/database"
	"github.com/IdrisAkintobi/go-basic-crud/routes"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file in local development if it exists
	appENV := strings.ToLower(os.Getenv("APP_ENV"))
	err := godotenv.Load()
	if appENV == "local" && err != nil {
		log.Fatalf("Warning: .env file not found or could not be loaded: %s", err)
	}
}

func gracefulShutdown(db *pgxpool.Pool) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	// Wait for signal
	<-ch
	signal.Stop(ch)

	// Do all necessary cleanup
	database.DisconnectDB(context.Background(), db)

	// Exit process
	os.Exit(0)
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect database
	conn, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	// Activate graceful shutdown
	go gracefulShutdown(conn)

	// Setup routes
	r := routes.SetupRoutes(conn)

	// Start server
	fmt.Printf("Server starting on %v\n", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}
