package services

import (
	"encoding/base64"
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
	sr                           *repository.SessionRepository
	sessionDuration, tokenLength int
}

func NewSessionService(db *pgx.Conn) *SessionService {
	sd, err := strconv.Atoi(os.Getenv("SESSION_DURATION"))
	if err != nil {
		panic("Invalid SESSION_DURATION: " + err.Error())
	}

	tl, err := strconv.Atoi(os.Getenv("TOKEN_LENGTH"))
	if err != nil {
		panic("Invalid TOKEN_LENGTH: " + err.Error())
	}

	return &SessionService{
		sr:              repository.NewSessionRepository(db),
		sessionDuration: sd,
		tokenLength:     tl,
	}
}

func (ss *SessionService) CreateSession(userId, deviceId, userAgent, ipAddress string) (token string, err error) {
	randomBytes, err := utils.RandomByte(uint32(ss.tokenLength))
	if err != nil {
		return "", fmt.Errorf("internal server error: %w", err)
	}

	token = base64.URLEncoding.EncodeToString(randomBytes)
	newSesParams := &schema.NewSessionParams{
		UserId:    userId,
		DeviceId:  deviceId,
		Token:     token,
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
	return session.Token, nil
}

func (ss *SessionService) UpdateSession(token string, duration uint) error {
	tokenHash := utils.Hash(token)
	expiresAt := time.Now().Add(time.Duration(duration) * time.Minute)

	err := ss.sr.ExtendSession(tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	return nil
}

func (ss *SessionService) DeleteSession(token string) error {
	tokenHash := utils.Hash(token)

	err := ss.sr.DeleteSession(tokenHash)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
