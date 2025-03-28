package repository

import (
	"context"

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
	INSERT INTO users (email, dob, firstName, lastName, passwordHash, createdAt, updatedAt)
	values ($1, $2, $3, $4, $5, $6, $7) returning *
	`, userData.Email, userData.DOB, userData.FirstName, userData.LastName, userData.PasswordHash, userData.CreatedAt, userData.UpdatedAt)

	err := row.Scan(&result.ID, &result.Email, &result.DOB, &result.FirstName, &result.LastName, &result.PasswordHash, &result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*schema.User, error) {
	var result schema.User
	row := r.db.QueryRow(context.Background(), `SELECT * FROM users WHERE email = $1`, email)

	err := row.Scan(&result.ID, &result.Email, &result.DOB, &result.FirstName, &result.LastName, &result.PasswordHash, &result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *UserRepository) GetUserById(userId string) (*schema.User, error) {
	var result schema.User
	row := r.db.QueryRow(context.Background(), `SELECT * FROM users WHERE id = $1`, userId)

	err := row.Scan(&result.ID, &result.Email, &result.DOB, &result.FirstName, &result.LastName, &result.PasswordHash, &result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
