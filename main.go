package main

import (
	"fmt"
	"log"

	"github.com/IdrisAkintobi/go-basic-crud/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	conn, err := database.ConnectDB()

	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	fmt.Println(conn)
}
