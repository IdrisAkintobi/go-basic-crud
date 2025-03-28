package services

import (
	"fmt"
	"time"

	"github.com/IdrisAkintobi/go-basic-crud/database/repository"
	"github.com/IdrisAkintobi/go-basic-crud/database/schema"
	"github.com/IdrisAkintobi/go-basic-crud/handlers/dto"
	"github.com/IdrisAkintobi/go-basic-crud/utils"
	"github.com/jackc/pgx/v5"
)

type UserService struct {
	ur *repository.UserRepository
}

func NewUserService(db *pgx.Conn) *UserService {
	return &UserService{ur: repository.NewUserRepository(db)}
}

func (us *UserService) RegisterUser(usr *dto.RegisterUserReqDTO) (*dto.RegisterUserResDTO, error) {

	// Parse user DOB string to time format
	dob, err := time.Parse(utils.DATE_LAYOUT, usr.DOB)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	// Hash password
	passwordHash, err := utils.Argon2id.GenerateHash([]byte(usr.Password), nil)
	if err != nil {
		return nil, fmt.Errorf("internal server error: %w", err)
	}

	newUserData := schema.NewUser(usr.Email, usr.FirstName, usr.LastName, passwordHash, dob)

	newUser, err := us.ur.CreateUser(newUserData)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Convert to the returned type
	data := &dto.RegisterUserResDTO{
		Email:     newUser.Email,
		DOB:       newUser.DOB.Format(utils.DATE_LAYOUT),
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		CreatedAt: newUser.CreatedAt,
	}

	return data, nil
}
