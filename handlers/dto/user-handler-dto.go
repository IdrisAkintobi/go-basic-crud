package dto

import "time"

type RegisterUserReqDTO struct {
	Email     string `json:"email"`
	DOB       string `json:"dob"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type RegisterUserResDTO struct {
	Email     string    `json:"email"`
	DOB       string    `json:"dob"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
}
