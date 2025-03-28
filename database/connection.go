package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type DbConfig struct {
	host, port, username, password, dbName string
}

func (config *DbConfig) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.username, config.password, config.host, config.port, config.dbName)
}

func ConnectDB() (*pgx.Conn, error) {
	dbCong := &DbConfig{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		username: os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
		dbName:   os.Getenv("DB_NAME"),
	}

	return pgx.Connect(context.Background(), dbCong.String())
}
