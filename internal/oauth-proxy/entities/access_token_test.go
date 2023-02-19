package entities

import (
	"testing"
	"time"
)

func TestIsExpired_ExpiredAccessToken(t *testing.T) {
	expiredAt := time.Now().Add(-24 * time.Hour)
	at := AccessToken{ExpiresAt: expiredAt}
	if !at.IsExpired() {
		t.Errorf("Expected IsExpired to return true for expired AccessToken")
	}
}

func TestIsExpired_NotExpiredAccessToken(t *testing.T) {
	notExpiredAt := time.Now().Add(24 * time.Hour)
	at := AccessToken{ExpiresAt: notExpiredAt}
	if at.IsExpired() {
		t.Errorf("Expected IsExpired to return false for not expired AccessToken")
	}
}
