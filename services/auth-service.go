package services

import (
	"errors"
	"fmt"

	"github.com/IdrisAkintobi/go-basic-crud/database/repository"
	"github.com/IdrisAkintobi/go-basic-crud/handlers/dto"
	"github.com/IdrisAkintobi/go-basic-crud/utils"
	"github.com/jackc/pgx/v5"
)

var ErrInvalidCred = errors.New("invalid login credentials")

type AuthService struct {
	ur   *repository.UserRepository
	ss   *SessionService
	g2ip *Geo2IPService
}

func NewAuthService(db *pgx.Conn) *AuthService {
	return &AuthService{
		ur:   repository.NewUserRepository(db),
		ss:   NewSessionService(db),
		g2ip: NewGeo2IPService(),
	}
}

func (as *AuthService) SignIn(cred *dto.AuthLoginReqDTO) (token string, err error) {

	user, err := as.ur.GetUserByEmail(cred.Email)
	if err != nil {
		return "", ErrInvalidCred
	}

	passwordMatch, err := utils.Argon2id.Compare(user.PasswordHash, []byte(cred.Password))
	if err != nil {
		return "", fmt.Errorf("internal server error: %w", err)
	}

	if !passwordMatch {
		return "", ErrInvalidCred
	}

	return as.ss.CreateSession(user.ID, cred.DeviceId, cred.UserAgent, cred.IPAddress)
}

func (as *AuthService) LogOut(tokenId int) error {
	err := as.ss.DeleteSessionById(tokenId)
	if err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}
	return nil
}

func (as *AuthService) WhoAmI(userId string) (*dto.WhoAmIResDTO, error) {
	user, err := as.ur.GetUserById(userId)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return nil, nil
	}

	userData := &dto.WhoAmIResDTO{
		ID:        user.ID,
		Email:     user.Email,
		DOB:       user.DOB.Format(utils.DATE_LAYOUT),
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return userData, nil

}

func (as *AuthService) GetActiveSessions(userId string) ([]*dto.ActiveSessionsDTO, error) {
	allSession, err := as.ss.FindAllActiveSession(userId)
	if err != nil {
		return nil, fmt.Errorf("error getting all user session: \n%w", err)
	}

	if allSession == nil {
		return []*dto.ActiveSessionsDTO{}, nil
	}

	var result []*dto.ActiveSessionsDTO

	for _, session := range allSession {
		location := as.g2ip.GetIPLocation(session.IPAddress)
		result = append(result, &dto.ActiveSessionsDTO{
			ID:        session.ID,
			UserAgent: session.UserAgent,
			IPAddress: session.IPAddress,
			Location:  location,
			CreatedAt: session.CreatedAt,
		})
	}

	return result, nil

}
