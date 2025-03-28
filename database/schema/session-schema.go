package schema

import "time"

type Session struct {
	ID        int       `db:"id"`
	UserId    string    `db:"userId"`
	Token     string    `db:"token"`
	UserAgent string    `db:"userAgent"`
	IPAddress string    `db:"ipAddress"`
	CreatedAt time.Time `db:"createdAt"`
	ExpiresAt time.Time `db:"expiresAt"`
}

// Constructor function to create a new User with default timestamps
func NewSession(userId, token, userAgent, ipAddress string, duration uint) *Session {
	now := time.Now()
	return &Session{
		UserId:    userId,
		Token:     token,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		CreatedAt: now,
		ExpiresAt: now.Add(time.Duration(duration) * time.Minute),
	}
}
