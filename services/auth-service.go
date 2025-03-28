package services

import (
	"errors"
	"fmt"

	"github.com/IdrisAkintobi/go-basic-crud/database/repository"
	"github.com/IdrisAkintobi/go-basic-crud/handlers/dto"
	"github.com/IdrisAkintobi/go-basic-crud/utils"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	ur *repository.UserRepository
	ss *SessionService
}

func NewAuthService(db *pgx.Conn) *AuthService {
	return &AuthService{
		ur: repository.NewUserRepository(db),
		ss: NewSessionService(db),
	}
}

var ErrInvalidCred = errors.New("invalid login credentials")

func (as *AuthService) SignIn(cred *dto.AuthLoginReqDTO) (token string, err error) {

	user, err := as.ur.GetUserByEmail(cred.Email)
	if err != nil {
		return "", ErrInvalidCred
	}

	passwordMatch, err := utils.Argon2id.Compare(user.PasswordHash, []byte(cred.Password))
	if err != nil {
		return "", fmt.Errorf("internal server error: %v", err)
	}

	if !passwordMatch {
		return "", ErrInvalidCred
	}

	token, err = as.ss.CreateSession(user.ID, cred.UserAgent, cred.IPAddress)
	return
}
