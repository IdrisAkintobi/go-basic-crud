package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/IdrisAkintobi/go-basic-crud/database/repository"
	"github.com/IdrisAkintobi/go-basic-crud/database/schema"
	"github.com/IdrisAkintobi/go-basic-crud/utils"
	"github.com/jackc/pgx/v5"
)

type SessionService struct {
	sr                                       *repository.SessionRepository
	sessionDuration, tokenLength, maxSession int
	SessionRefreshWindow                     time.Duration
}

func NewSessionService(db *pgx.Conn) *SessionService {
	sd, err := strconv.Atoi(os.Getenv("SESSION_DURATION"))
	if err != nil {
		panic("Invalid SESSION_DURATION: " + err.Error())
	}

	srw, err := strconv.Atoi(os.Getenv("SESSION_REFRESH_WINDOW"))
	if err != nil {
		panic("Invalid SESSION_REFRESH_WINDOW: " + err.Error())
	}

	tl, err := strconv.Atoi(os.Getenv("TOKEN_LENGTH"))
	if err != nil {
		panic("Invalid TOKEN_LENGTH: " + err.Error())
	}

	mxS, err := strconv.Atoi(os.Getenv("MAXIMUM_SESSION"))
	if err != nil {
		panic("Invalid MAXIMUM_SESSION: " + err.Error())
	}

	return &SessionService{
		sr:                   repository.NewSessionRepository(db),
		sessionDuration:      sd,
		SessionRefreshWindow: time.Minute * time.Duration(srw),
		tokenLength:          tl,
		maxSession:           mxS,
	}
}

var ErrMaximumSession = errors.New("maximum session reached")
var ErrInternalServer = errors.New("internal server error")

func (ss *SessionService) CreateSession(userId, deviceId, userAgent, ipAddress string) (token string, err error) {
	// Check if maximum session is reached
	sessionCount, err := ss.sr.CountUserActiveSessions(userId)
	if err != nil {
		return "", fmt.Errorf("internal server error: %w", err)
	}
	if sessionCount >= ss.maxSession {
		return "", ErrMaximumSession
	}

	// Generate token
	randomBytes, err := utils.RandomByte(uint32(ss.tokenLength))
	if err != nil {
		return "", fmt.Errorf("internal server error: %w", err)
	}
	token = base64.URLEncoding.EncodeToString(randomBytes)

	// Hash session token to be saved to the db
	tokenHash := utils.Hash(token)

	newSesParams := &schema.NewSessionParams{
		UserId:    userId,
		DeviceId:  deviceId,
		Token:     tokenHash,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		Duration:  uint(ss.sessionDuration),
	}
	session := schema.NewSession(newSesParams)

	_, err = ss.sr.CreateSession(session)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	// Return the actual token
	return token, nil
}

func (ss *SessionService) FindSession(token string) (*schema.Session, error) {
	tokenHash := utils.Hash(token)

	session, err := ss.sr.FindSession(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	return session, err

}

func (ss *SessionService) FindAllActiveSession(userId string) ([]*schema.Session, error) {
	return ss.sr.FindActiveSession(userId)
}

func (ss *SessionService) ExtendSession(token string) error {
	tokenHash := utils.Hash(token)
	expiresAt := time.Now().Add(time.Duration(ss.sessionDuration) * time.Minute)

	err := ss.sr.ExtendSession(tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	return nil
}

func (ss *SessionService) DeleteSessionByToken(token string) error {
	tokenHash := utils.Hash(token)

	err := ss.sr.DeleteSessionByToken(tokenHash)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

func (ss *SessionService) DeleteSessionById(id int) error {
	err := ss.sr.DeleteSessionById(id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

func (ss *SessionService) ClearUserSession(userId string) error {
	err := ss.sr.ClearUserSession(userId)
	if err != nil {
		return fmt.Errorf("failed to clear user session: %w", err)
	}
	return nil
}
