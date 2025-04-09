package repository_test

import (
	"context"
	"fmt"
	"time"

	"github.com/IdrisAkintobi/go-basic-crud/database/repository"
	"github.com/IdrisAkintobi/go-basic-crud/database/schema"
	"github.com/IdrisAkintobi/go-basic-crud/utils"
	"github.com/jackc/pgx/v5"
)

const mockUUID = "f4bd11c8-840a-45eb-b72a-917bd90996a7"

var mockUser = &schema.User{
	Email:        "john.doe@example.com",
	DOB:          time.Date(1990, 1, 11, 0, 0, 0, 0, time.UTC),
	FirstName:    "John",
	LastName:     "Doe",
	PasswordHash: []byte(""),
}

func countUsers(db *pgx.Conn) (int, error) {
	var count int
	err := db.QueryRow(context.Background(), `
	SELECT count(*) FROM users;
	`).Scan(&count)

	return count, err
}

func (ts *RepositoryTestSuite) TestCreateUser() {
	// Count users in db before creating user
	before, err := countUsers(ts.db)
	ts.NoError(err)

	// Create user repository
	ur := repository.NewUserRepository(ts.db)

	// Create user
	dbUser, err := ur.CreateUser(mockUser)
	ts.NoError(err)

	// Count users in db after creating user
	after, err := countUsers(ts.db)
	ts.NoError(err)

	// Assert
	ts.Greater(after, before)
	ts.Equal(after, before+1)
	ts.True(dbUser.CreatedAt.Equal(dbUser.UpdatedAt))
	ts.Equal(mockUser.Email, dbUser.Email)
}

func (ts *RepositoryTestSuite) TestPasswordHash() {
	// Create user repository
	ur := repository.NewUserRepository(ts.db)

	// Update password
	passwordHash, err := utils.Argon2id.GenerateHash([]byte("password123"), nil)
	ts.NoError(err)
	mockUser.PasswordHash = passwordHash

	// Save user to the db
	_, err = ur.CreateUser(mockUser)
	ts.NoError(err)

	// Get user
	dbUser, err := ur.GetUserByEmail(mockUser.Email)
	ts.NoError(err)

	// Assert
	ts.Equal(mockUser.Email, dbUser.Email)
	ts.Equal(passwordHash, dbUser.PasswordHash)
}

func (ts *RepositoryTestSuite) TestGetUserByEmail() {
	ts.TestCreateUser()
	// Create user repository
	ur := repository.NewUserRepository(ts.db)

	// Get user
	dbUser, err := ur.GetUserByEmail(mockUser.Email)
	ts.NoError(err)

	// Assert
	ts.Equal(mockUser.Email, dbUser.Email)
	ts.Equal(mockUser.PasswordHash, dbUser.PasswordHash)

	//Get non-existing user
	dbUser, err = ur.GetUserByEmail(fmt.Sprintf("non-existing-%s", mockUser.Email))
	ts.Nil(err)
	ts.Nil(dbUser)
}

func (ts *RepositoryTestSuite) TestGetUserById() {
	ts.TestCreateUser()
	// Create user repository
	ur := repository.NewUserRepository(ts.db)

	// Get user
	dbUser, err := ur.GetUserByEmail(mockUser.Email)
	ts.NoError(err)
	dbUser, err = ur.GetUserById(dbUser.ID)
	ts.NoError(err)

	// Assert
	ts.Equal(mockUser.Email, dbUser.Email)
	ts.Equal(mockUser.PasswordHash, dbUser.PasswordHash)

	//Get non-existing user
	dbUser, err = ur.GetUserById(mockUUID)
	ts.Nil(dbUser)
	ts.Nil(err)
}
