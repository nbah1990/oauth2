package entities

import (
	"time"
)

type RefreshToken struct {
	ID            string
	AccessTokenID string
	Revoked       bool
	ExpiresAt     time.Time
}

func (at RefreshToken) IsExpired() bool {
	return time.Now().After(at.ExpiresAt)
}
