package repository_test

import (
	"context"
	"testing"

	"github.com/IdrisAkintobi/go-basic-crud/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	db *pgxpool.Pool
}

func (ts *RepositoryTestSuite) SetupTest() {
	cfg := config.Get()
	conn, err := pgxpool.New(context.Background(), cfg.TestDatabaseURL)
	if err != nil {
		panic((err))
	}
	conn.Exec(context.Background(), `TRUNCATE sessions`)
	conn.Exec(context.Background(), `TRUNCATE users CASCADE`)
	ts.db = conn
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (ts *RepositoryTestSuite) TearDownSuite() {
	// Close the database connection
	ts.db.Close()
}
