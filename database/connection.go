package database

import (
	"context"
	"fmt"
	"time"

	"github.com/IdrisAkintobi/go-basic-crud/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConfig struct {
	host, port, username, password, dbName string
}

func (config *DbConfig) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.username, config.password, config.host, config.port, config.dbName)
}

func ConnectDB() (*pgxpool.Pool, error) {
	cfg := config.Get()
	dbCong := &DbConfig{
		host:     cfg.DBHost,
		port:     cfg.DBPort,
		username: cfg.DBUser,
		password: cfg.DBPassword,
		dbName:   cfg.DBName,
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

	done := make(chan struct{})

	go func() {
		db.Close()
		close(done)
	}()

	select {
	case <-ctx.Done():
		err := fmt.Errorf("\nDatabase graceful shutdown canceled: %w", ctx.Err())
		fmt.Println(err)
	case <-done:
		fmt.Println("\nDatabase gracefully shutdown")
	}
}
