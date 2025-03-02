package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	db *pgx.Conn
}

func (ts *RepositoryTestSuite) SetupTest() {
	dbConnStr := os.Getenv("TEST_DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), dbConnStr)
	if err != nil {
		panic((err))
	}
	conn.Exec(context.Background(), `TRUNCATE users CASCADE`)
	ts.db = conn
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (ts *RepositoryTestSuite) TearDownSuite() {
	// Close the database connection
	ts.db.Close(context.Background())
}
