package repository

import (
	"context"
	"errors"

	"github.com/IdrisAkintobi/go-basic-crud/database/schema"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(userData *schema.User) (*schema.User, error) {
	var result schema.User
	row := r.db.QueryRow(context.Background(), `
	INSERT INTO users (email, dob, firstName, lastName, PasswordHash, createdAt, updatedAt)
	values ($1, $2, $3, $4, $5, $6, $7) returning id, email, dob, firstName, lastName, passwordHash
	`, userData.Email, userData.DOB, userData.FirstName, userData.LastName, userData.PasswordHash, userData.CreatedAt, userData.UpdatedAt)

	err := row.Scan(&result.ID, &result.Email, &result.DOB, &result.FirstName, &result.LastName, &result.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*schema.User, error) {
	var result schema.User
	row := r.db.QueryRow(context.Background(), `
	SELECT id, email, dob, firstName, lastName, PasswordHash
	FROM users
	WHERE email = $1`, email)

	err := row.Scan(&result.ID, &result.Email, &result.DOB, &result.FirstName, &result.LastName, &result.PasswordHash)

	if err != nil {
		return handleFindUserError(err)
	}

	return &result, nil
}

func (r *UserRepository) GetUserById(userId string) (*schema.User, error) {
	var result schema.User
	row := r.db.QueryRow(context.Background(), `
	SELECT id, email, dob, firstName, lastName, PasswordHash
	FROM users WHERE id = $1`, userId)

	err := row.Scan(&result.ID, &result.Email, &result.DOB, &result.FirstName, &result.LastName, &result.PasswordHash)

	if err != nil {
		return handleFindUserError(err)
	}

	return &result, nil
}

func handleFindUserError(err error) (*schema.User, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return nil, err
}
