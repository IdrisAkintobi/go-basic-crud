package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IdrisAkintobi/go-basic-crud/database/schema"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(sessionData *schema.Session) (*schema.Session, error) {
	row := r.db.QueryRow(context.Background(), `
	INSERT INTO sessions (userId, deviceId, token, userAgent, ipAddress, createdAt, expiresAt)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (userId, deviceId) 
	DO UPDATE SET 
		token = EXCLUDED.token, 
		userAgent = EXCLUDED.userAgent, 
		ipAddress = EXCLUDED.ipAddress, 
		createdAt = EXCLUDED.createdAt, 
		expiresAt = EXCLUDED.expiresAt
	RETURNING id
	`, sessionData.UserId, sessionData.DeviceId, sessionData.Token, sessionData.UserAgent, sessionData.IPAddress, sessionData.CreatedAt, sessionData.ExpiresAt)

	err := row.Scan(&sessionData.ID)
	if err != nil {
		return nil, err
	}

	return sessionData, nil
}

func (r *SessionRepository) FindSession(tokenHash string) (*schema.Session, error) {
	var result schema.Session

	row := r.db.QueryRow(context.Background(), `
	SELECT id, userId, deviceId, token, userAgent, ipAddress, expiresAt
	FROM sessions 
	WHERE token = $1`, tokenHash)

	err := row.Scan(&result.ID, &result.UserId, &result.DeviceId, &result.Token, &result.UserAgent, &result.IPAddress, &result.ExpiresAt)
	if err != nil {
		return handleFindSessionError(err)
	}

	return &result, nil
}

func (r *SessionRepository) DeleteExistingDeviceSession(userId, deviceId string) error {
	_, err := r.db.Exec(context.Background(), `
	DELETE FROM sessions 
	WHERE userId = $1 AND deviceId = $2`, userId, deviceId)
	return err
}

func (r *SessionRepository) FindActiveSession(userId string) ([]*schema.Session, error) {
	var sessions []*schema.Session

	rows, err := r.db.Query(context.Background(), `
	SELECT id, userId, deviceId, token, userAgent, ipAddress, createdAt, expiresAt
	FROM sessions
	WHERE userId = $1 AND expiresAt > NOW()`, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query for sessions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var session schema.Session
		err = rows.Scan(&session.ID, &session.UserId, &session.DeviceId, &session.Token, &session.UserAgent, &session.IPAddress, &session.CreatedAt, &session.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session row: %w", err)
		}

		sessions = append(sessions, &session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate sessions: %w", err)
	}

	return sessions, nil
}

func (r *SessionRepository) ExtendSession(tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(context.Background(), `UPDATE sessions SET expiresAt = $1 WHERE token = $2`, expiresAt, tokenHash)
	return err
}

func (r *SessionRepository) CountUserActiveSessions(userId string) (int, error) {
	var count int
	row := r.db.QueryRow(context.Background(), `SELECT COUNT(*) FROM sessions WHERE userId = $1 AND expiresAt > NOW()`, userId)
	err := row.Scan(&count)
	return count, err
}

func (r *SessionRepository) DeleteSessionById(id int) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM sessions WHERE id = $1`, id)
	return err
}

func (r *SessionRepository) DeleteSessionByToken(tokenHash string) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM sessions WHERE token = $1`, tokenHash)
	return err
}

func (r *SessionRepository) ClearUserSession(userId string) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM sessions WHERE userId = $1`, userId)
	return err
}

func handleFindSessionError(err error) (*schema.Session, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return nil, err
}
