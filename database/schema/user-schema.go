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
