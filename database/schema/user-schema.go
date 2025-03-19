package schema

import "time"

type User struct {
	ID           int       `db:"id"`
	Email        string    `db:"email"`
	DOB          time.Time `db:"dob"`
	FirstName    string    `db:"firstName"`
	LastName     string    `db:"lastName"`
	PasswordHash string    `db:"passwordHash"`
	CreatedAt    time.Time `db:"createdAt"`
	UpdatedAt    time.Time `db:"updatedAt"`
}

// Constructor function to create a new User with default timestamps
func NewUser(email, firstName, lastName, passwordHash string, dob time.Time) *User {
	now := time.Now()
	return &User{
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
