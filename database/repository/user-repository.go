package repository

import (
	"context"

	"github.com/IdrisAkintobi/go-basic-crud/database/schema"
	"github.com/IdrisAkintobi/go-basic-crud/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(userData *schema.User) (*schema.User, error) {
	row := r.db.QueryRow(context.Background(), `
	INSERT INTO users (email, dob, firstName, lastName, PasswordHash, createdAt, updatedAt)
	values ($1, $2, $3, $4, $5, $6, $7) returning id
	`, userData.Email, userData.DOB, userData.FirstName, userData.LastName, userData.PasswordHash, userData.CreatedAt, userData.UpdatedAt)

	err := row.Scan(&userData.ID)

	return userData, utils.FormatDBError(err)
}

func (r *UserRepository) GetUserByEmail(email string) (*schema.User, error) {
	var result schema.User
	row := r.db.QueryRow(context.Background(), `
	SELECT id, email, dob, firstName, lastName, PasswordHash
	FROM users
	WHERE email = $1`, email)

	err := row.Scan(&result.ID, &result.Email, &result.DOB, &result.FirstName, &result.LastName, &result.PasswordHash)

	return &result, utils.FormatDBError(err)
}

func (r *UserRepository) GetUserById(userId string) (*schema.User, error) {
	var result schema.User
	row := r.db.QueryRow(context.Background(), `
	SELECT id, email, dob, firstName, lastName, PasswordHash
	FROM users WHERE id = $1`, userId)

	err := row.Scan(&result.ID, &result.Email, &result.DOB, &result.FirstName, &result.LastName, &result.PasswordHash)

	return &result, utils.FormatDBError(err)
}
