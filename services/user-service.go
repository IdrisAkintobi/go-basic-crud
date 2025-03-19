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

	dob, err := time.Parse(utils.DATE_LAYOUT, usr.DOB)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	newUserData := schema.NewUser(usr.Email, usr.FirstName, usr.LastName, usr.Password, dob)

	newUser, err := us.ur.CreateUser(newUserData)
	if err != nil {
		return nil, err
	}
	// Convert to the returned type
	data := &dto.RegisterUserResDTO{
		Email:     newUser.Email,
		DOB:       newUser.DOB.Local().Format(utils.DATE_LAYOUT),
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		CreatedAt: newUser.CreatedAt,
	}

	return data, nil
}
