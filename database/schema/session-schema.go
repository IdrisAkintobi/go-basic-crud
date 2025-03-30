package schema

import "time"

type Session struct {
	ID        int       `db:"id"`
	UserId    string    `db:"userId"`
	DeviceId  string    `db:"deviceId"`
	Token     string    `db:"token"`
	UserAgent string    `db:"userAgent"`
	IPAddress string    `db:"ipAddress"`
	CreatedAt time.Time `db:"createdAt"`
	ExpiresAt time.Time `db:"expiresAt"`
}

type NewSessionParams struct {
	UserId, DeviceId, Token, UserAgent, IPAddress string
	Duration                                      uint
}

// Constructor function to create a new User with default timestamps
func NewSession(params *NewSessionParams) *Session {
	now := time.Now()
	return &Session{
		UserId:    params.UserId,
		DeviceId:  params.DeviceId,
		Token:     params.Token,
		UserAgent: params.UserAgent,
		IPAddress: params.IPAddress,
		CreatedAt: now,
		ExpiresAt: now.Add(time.Duration(params.Duration) * time.Minute),
	}
}
