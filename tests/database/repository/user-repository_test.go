package repository_test

import (
	"context"
	"fmt"
	"time"

	"github.com/IdrisAkintobi/go-basic-crud/database/repository"
	"github.com/IdrisAkintobi/go-basic-crud/database/schema"
	"github.com/jackc/pgx/v5"
)

var testUser = &schema.User{
	Email:        "john.doe@example.com",
	DOB:          time.Date(1991, 4, 17, 0, 0, 0, 0, time.UTC),
	FirstName:    "John",
	LastName:     "Doe",
	PasswordHash: "password123",
}

func countDB(db *pgx.Conn) (int, error) {
	var count int
	err := db.QueryRow(context.Background(), `
	SELECT count(*) FROM users;
	`).Scan(&count)

	return count, err
}

func (ts *RepositoryTestSuite) TestInsertUser() {
	// Count users in db before creating user
	before, err := countDB(ts.db)
	ts.NoError(err)

	// Create user repository
	ur := repository.NewUserRepository(ts.db)

	// Create user
	dbUser, err := ur.CreateUser(testUser)
	ts.NoError(err)

	// Count users in db after creating user
	after, err := countDB(ts.db)
	ts.NoError(err)

	// Assert
	ts.Greater(dbUser.ID, 0)
	ts.Greater(after, before)
	ts.Equal(after, before+1)
	ts.True(dbUser.CreatedAt.Equal(dbUser.UpdatedAt))
	ts.Equal(testUser.Email, dbUser.Email)
}

func (ts *RepositoryTestSuite) TestGetUserByEmail() {
	ts.TestInsertUser()
	// Create user repository
	ur := repository.NewUserRepository(ts.db)

	// Get user
	dbUser, err := ur.GetUserByEmail(testUser.Email)
	ts.NoError(err)

	// Assert
	ts.Equal(testUser.Email, dbUser.Email)
	ts.Equal(testUser.PasswordHash, dbUser.PasswordHash)

	//Get non-existing user
	dbUser, err = ur.GetUserByEmail(fmt.Sprintf("non-existing-%s", testUser.Email))
	ts.Error(err)
	ts.ErrorIs(err, pgx.ErrNoRows)
	ts.Nil(dbUser)
}

func (ts *RepositoryTestSuite) TestGetUserById() {
	ts.TestInsertUser()
	// Create user repository
	ur := repository.NewUserRepository(ts.db)

	// Get user
	dbUser, err := ur.GetUserByEmail(testUser.Email)
	ts.NoError(err)
	dbUser, err = ur.GetUserById(dbUser.ID)
	ts.NoError(err)

	// Assert
	ts.Equal(testUser.Email, dbUser.Email)
	ts.Equal(testUser.PasswordHash, dbUser.PasswordHash)

	//Get non-existing user
	dbUser, err = ur.GetUserByEmail(fmt.Sprintf("%d", testUser.ID+12))
	ts.Error(err)
	ts.ErrorIs(err, pgx.ErrNoRows)
	ts.Nil(dbUser)
}
