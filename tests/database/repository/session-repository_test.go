package repository_test

import (
	"context"
	"time"

	"github.com/IdrisAkintobi/go-basic-crud/database/repository"
	"github.com/IdrisAkintobi/go-basic-crud/database/schema"
	"github.com/IdrisAkintobi/go-basic-crud/utils"
	"github.com/jackc/pgx/v5"
)

const (
	userId    string = "99533d94-b3d7-43a7-972d-d91d81911033"
	deviceId  string = "4cc10adf-320c-455c-95c3-14830c18676d"
	token     string = "base64Token=="
	userAgent string = "Go Test"
	ipAddress string = "127.0.0.1"
	duration  uint   = 60
)

var mockSession = schema.NewSession(&schema.NewSessionParams{
	UserId:    userId,
	DeviceId:  deviceId,
	Token:     token,
	UserAgent: userAgent,
	IPAddress: ipAddress,
	Duration:  duration,
})

func countSessions(db *pgx.Conn) (int, error) {
	var count int
	err := db.QueryRow(context.Background(), `
	SELECT count(*) FROM sessions;
	`).Scan(&count)

	return count, err
}

func (ts *RepositoryTestSuite) TestCreateSession() {
	// Count session in db before creating session
	before, err := countSessions(ts.db)
	ts.NoError(err)

	// Create session repository
	sr := repository.NewSessionRepository(ts.db)

	// Create session
	dbSession, err := sr.CreateSession(mockSession)
	ts.NoError(err)

	// Count session in db after creating session
	after, err := countSessions(ts.db)
	ts.NoError(err)

	// Assert
	ts.Greater(after, before)
	ts.Equal(after, before+1)
	ts.Equal(dbSession.UserAgent, mockSession.UserAgent)
	// Ensure token is hashed
	ts.NotEqual(dbSession.Token, mockSession.Token)
	ts.Equal(dbSession.Token, utils.Hash(mockSession.Token))
}

func (ts *RepositoryTestSuite) TestFindSession() {
	// Create session repository
	sr := repository.NewSessionRepository(ts.db)

	// Create session
	_, err := sr.CreateSession(mockSession)
	ts.NoError(err)

	dbSession, err := sr.FindSession(mockSession.Token)
	ts.NoError(err)

	// Assert
	ts.Equal(dbSession.UserId, mockSession.UserId)
	ts.Equal(dbSession.UserAgent, mockSession.UserAgent)
	ts.Equal(dbSession.IPAddress, mockSession.IPAddress)
	ts.Equal(dbSession.DeviceId, mockSession.DeviceId)
}

func (ts *RepositoryTestSuite) TestUpdateSession() {
	// Create session repository
	sr := repository.NewSessionRepository(ts.db)

	// Create session
	dbSession, err := sr.CreateSession(mockSession)
	ts.NoError(err)

	duration := time.Now().Add(time.Hour)
	// Update session
	err = sr.ExtendSession(mockSession.Token, duration)
	ts.NoError(err)

	updatedSession, err := sr.FindSession(mockSession.Token)
	ts.NoError(err)

	// Assert
	ts.Equal(dbSession.UserAgent, updatedSession.UserAgent)
	ts.Equal(dbSession.Token, updatedSession.Token)
	ts.Equal(dbSession.IPAddress, updatedSession.IPAddress)
	ts.NotEqual(dbSession.ExpiresAt, updatedSession.ExpiresAt)
	ts.Equal(updatedSession.ExpiresAt.Compare(dbSession.ExpiresAt), +1)
}

func (ts *RepositoryTestSuite) TestDeleteSession() {
	// Count session in db before creating session
	beforeCreate, err := countSessions(ts.db)
	ts.NoError(err)

	// Create session repository
	sr := repository.NewSessionRepository(ts.db)

	// Create session
	_, err = sr.CreateSession(mockSession)
	ts.NoError(err)

	// Count session in db after creating session
	afterCreate, err := countSessions(ts.db)
	ts.NoError(err)

	// Delete session
	err = sr.DeleteSession(mockSession.Token)
	ts.NoError(err)

	// Count session in db after creating session
	afterDelete, err := countSessions(ts.db)
	ts.NoError(err)

	// Assert
	_, err = sr.FindSession(mockSession.Token)
	ts.Error(err)
	ts.Greater(afterCreate, beforeCreate)
	ts.Equal(beforeCreate, afterDelete)
}
