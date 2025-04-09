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
	deviceId2 string = "4cc10adf-320c-455c-95c3-14830c18676d"
	token     string = "base64Token=="
	token2    string = "base64Token2=="
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

var mockSession2 = schema.NewSession(&schema.NewSessionParams{
	UserId:    userId,
	DeviceId:  deviceId2,
	Token:     token2,
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

	tokenHash := utils.Hash(mockSession.Token)
	dbSession, err := sr.FindSession(tokenHash)
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
	tokenHash := utils.Hash(mockSession.Token)
	// Update session
	err = sr.ExtendSession(tokenHash, duration)
	ts.NoError(err)

	updatedSession, err := sr.FindSession(tokenHash)
	ts.NoError(err)

	// Assert
	ts.Equal(dbSession.UserAgent, updatedSession.UserAgent)
	ts.Equal(dbSession.Token, updatedSession.Token)
	ts.Equal(dbSession.IPAddress, updatedSession.IPAddress)
	ts.NotEqual(dbSession.ExpiresAt, updatedSession.ExpiresAt)
	ts.Equal(updatedSession.ExpiresAt.Compare(dbSession.ExpiresAt), +1)
}

func (ts *RepositoryTestSuite) TestCountActiveSession() {
	// Create session repository
	sr := repository.NewSessionRepository(ts.db)

	// Count active session
	beforeCreate, err := sr.CountUserActiveSessions(mockSession.UserId)
	ts.NoError(err)

	// Create session
	_, err = sr.CreateSession(mockSession)
	ts.NoError(err)

	// Count session in db after creating session
	firstCreate, err := sr.CountUserActiveSessions(mockSession.UserId)
	ts.NoError(err)

	// Create another session
	_, err = sr.CreateSession(mockSession2)
	ts.NoError(err)

	// Count session in db after creating session
	secondCreate, err := sr.CountUserActiveSessions(mockSession.UserId)
	ts.NoError(err)

	// Assert
	ts.Greater(firstCreate, beforeCreate)
	ts.Greater(secondCreate, firstCreate)
	ts.Equal(firstCreate, 1)
	ts.Equal(secondCreate, 2)
}

func (ts *RepositoryTestSuite) TestFindAllSession() {
	// Create session repository
	sr := repository.NewSessionRepository(ts.db)

	// Create session
	_, err := sr.CreateSession(mockSession)
	ts.NoError(err)

	// Create another session
	_, err = sr.CreateSession(mockSession2)
	ts.NoError(err)

	// Count session in db after creating session
	allSessions, err := sr.FindAllSession(mockSession.UserId)
	ts.NoError(err)

	// Assert
	ts.Equal(allSessions[0].UserId, mockSession.UserId)
	ts.Equal(allSessions[1].UserId, mockSession2.UserId)
}

func (ts *RepositoryTestSuite) TestDeleteSession() {
	// Count session in db before creating session
	beforeCreate, err := countSessions(ts.db)
	ts.NoError(err)

	// Create session repository
	sr := repository.NewSessionRepository(ts.db)

	// Create first session
	firstSession, err := sr.CreateSession(mockSession)
	ts.NoError(err)

	// Count session in db after first session creation
	firstCreateCount, err := countSessions(ts.db)
	ts.NoError(err)

	// Create second session
	_, err = sr.CreateSession(mockSession2)
	ts.NoError(err)

	// Count session in db after second session creation
	secondCreateCount, err := countSessions(ts.db)
	ts.NoError(err)

	// Delete first session by id
	err = sr.DeleteSessionById(firstSession.ID)
	ts.NoError(err)

	// Count session in db after deleting first session
	afterFirstDelete, err := countSessions(ts.db)
	ts.NoError(err)

	// Delete second session by token
	tokenHash := utils.Hash(mockSession2.Token)
	err = sr.DeleteSessionByToken(tokenHash)
	ts.NoError(err)

	// Count session in db after deleting first session
	afterSecondDelete, err := countSessions(ts.db)
	ts.NoError(err)

	// Assert
	dbSession, err := sr.FindSession(utils.Hash(mockSession.Token))
	ts.Nil(dbSession)
	ts.Nil(err)

	ts.Greater(firstCreateCount, beforeCreate)
	ts.Greater(secondCreateCount, firstCreateCount)
	ts.Equal(secondCreateCount, 2)
	ts.Less(afterSecondDelete, afterFirstDelete)
	ts.Equal(afterSecondDelete, beforeCreate)
}

func (ts *RepositoryTestSuite) TestClearUserSession() {
	// Count session in db before creating session
	beforeCreate, err := countSessions(ts.db)
	ts.NoError(err)

	// Create session repository
	sr := repository.NewSessionRepository(ts.db)

	// Create first session
	_, err = sr.CreateSession(mockSession)
	ts.NoError(err)

	// Count session in db after first session creation
	firstCreateCount, err := countSessions(ts.db)
	ts.NoError(err)

	// Create second session
	_, err = sr.CreateSession(mockSession2)
	ts.NoError(err)

	// Count session in db after second session creation
	secondCreateCount, err := countSessions(ts.db)
	ts.NoError(err)

	// Clear user session
	err = sr.ClearUserSession(mockSession.UserId)
	ts.NoError(err)

	// Count session after clearing
	afterClearing, err := countSessions(ts.db)
	ts.NoError(err)

	// Assert
	dbSession, err := sr.FindSession(utils.Hash(mockSession.Token))
	ts.Nil(dbSession)
	ts.Nil(err)

	ts.Greater(firstCreateCount, beforeCreate)
	ts.Greater(secondCreateCount, firstCreateCount)
	ts.Equal(secondCreateCount, 2)
	ts.Equal(beforeCreate, afterClearing)
}
