package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConfig struct {
	host, port, username, password, dbName string
}

func (config *DbConfig) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.username, config.password, config.host, config.port, config.dbName)
}

func ConnectDB() (*pgxpool.Pool, error) {
	dbCong := &DbConfig{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		username: os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
		dbName:   os.Getenv("DB_NAME"),
	}

	config, err := pgxpool.ParseConfig(dbCong.String())
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database:\n%w", err)
	}

	// Set minimum and maximum connection to 2 - Ensure we always have 2 connections in pool
	config.MinConns = 2
	config.MaxConns = 2

	return pgxpool.NewWithConfig(context.Background(), config)
}

func DisconnectDB(ctx context.Context, db *pgxpool.Pool) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		fmt.Printf("\nDatabase graceful shutdown canceled:\n%v", ctx.Err())
	default:
		db.Close()
		fmt.Println("\nDatabase gracefully shutdown")
	}
}
