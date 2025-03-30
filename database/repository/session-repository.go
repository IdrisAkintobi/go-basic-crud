package repository

import (
	"context"
	"time"

	"github.com/IdrisAkintobi/go-basic-crud/database/schema"
	"github.com/IdrisAkintobi/go-basic-crud/utils"
	"github.com/jackc/pgx/v5"
)

type SessionRepository struct {
	db *pgx.Conn
}

func NewSessionRepository(db *pgx.Conn) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(sessionData *schema.Session) (*schema.Session, error) {
	var result schema.Session

	// Hash session token before saving to the db
	tokenHash := utils.Hash(sessionData.Token)

	row := r.db.QueryRow(context.Background(), `
	INSERT INTO sessions (userId, deviceId, token, userAgent, ipAddress, createdAt, expiresAt)
	values ($1, $2, $3, $4, $5, $6, $7) returning *
	`, sessionData.UserId, sessionData.DeviceId, tokenHash, sessionData.UserAgent, sessionData.IPAddress, sessionData.CreatedAt, sessionData.ExpiresAt)

	err := row.Scan(&result.ID, &result.UserId, &result.DeviceId, &result.Token, &result.UserAgent, &result.IPAddress, &result.CreatedAt, &result.ExpiresAt)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *SessionRepository) FindSession(token string) (*schema.Session, error) {
	var result schema.Session
	tokenHash := utils.Hash(token)

	row := r.db.QueryRow(context.Background(), `SELECT * FROM sessions WHERE token = $1`, tokenHash)

	err := row.Scan(&result.ID, &result.UserId, &result.DeviceId, &result.Token, &result.UserAgent, &result.IPAddress, &result.CreatedAt, &result.ExpiresAt)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *SessionRepository) ExtendSession(token string, expiresAt time.Time) error {
	tokenHash := utils.Hash(token)
	_, err := r.db.Exec(context.Background(), `UPDATE sessions SET expiresAt = $1 WHERE token = $2`, expiresAt, tokenHash)

	return err
}

func (r *SessionRepository) DeleteSession(token string) error {
	tokenHash := utils.Hash(token)
	_, err := r.db.Exec(context.Background(), `DELETE FROM sessions WHERE token = $1`, tokenHash)

	return err
}
