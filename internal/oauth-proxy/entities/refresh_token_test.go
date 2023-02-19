package entities

import (
	"testing"
	"time"
)

func TestIsExpired_ExpiredRefreshToken(t *testing.T) {
	expiredAt := time.Now().Add(-24 * time.Hour)
	at := RefreshToken{ExpiresAt: expiredAt}
	if !at.IsExpired() {
		t.Errorf("Expected IsExpired to return true for expired RefreshToken")
	}
}

func TestIsExpired_NotExpiredRefreshToken(t *testing.T) {
	notExpiredAt := time.Now().Add(24 * time.Hour)
	at := RefreshToken{ExpiresAt: notExpiredAt}
	if at.IsExpired() {
		t.Errorf("Expected IsExpired to return false for not expired RefreshToken")
	}
}
