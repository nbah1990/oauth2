package entities

import (
	"time"
)

type AccessToken struct {
	ID        string
	UserID    string
	ClientID  string
	Revoked   bool
	ExpiresAt time.Time
}

func (at AccessToken) IsExpired() bool {
	return time.Now().After(at.ExpiresAt)
}
