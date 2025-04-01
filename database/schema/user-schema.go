package schema

import "time"

type User struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	DOB          time.Time `db:"dob"`
	FirstName    string    `db:"firstName"`
	LastName     string    `db:"lastName"`
	PasswordHash []byte    `db:"passwordHash"`
	CreatedAt    time.Time `db:"createdAt"`
	UpdatedAt    time.Time `db:"updatedAt"`
}

// Constructor function to create a new User with default timestamps
func NewUser(email, firstName, lastName string, passwordHash []byte, dob time.Time) *User {
	now := time.Now()
	return &User{
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: passwordHash,
		DOB:          dob,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
