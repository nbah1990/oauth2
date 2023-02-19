package entities

import (
	"oauth-proxy/internal/oauth-proxy/services/hash"
	"testing"
)

func TestUserCheckPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := hash.BcryptHash(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	user := &User{
		Username: "john_doe",
		Password: hashedPassword,
	}

	// Test with correct password
	if !user.CheckPassword(password) {
		t.Error("Expected user.CheckPassword to return true for correct password")
	}

	// Test with incorrect password
	if user.CheckPassword("incorrect_password") {
		t.Error("Expected user.CheckPassword to return false for incorrect password")
	}
}

func TestUserCheckPassword2(t *testing.T) {
	password := "password123"
	user := CreateUser("john_doe", password, "test_id")

	// Test with correct password
	if !user.CheckPassword(password) {
		t.Error("Expected user.CheckPassword to return true for correct password")
	}

	// Test with incorrect password
	if user.CheckPassword("incorrect_password") {
		t.Error("Expected user.CheckPassword to return false for incorrect password")
	}
}

func TestCreateUser(t *testing.T) {
	username := "john"
	password := "doe"

	u := CreateUser(username, password, "test_id")
	if u.Username != username {
		t.Fatalf("Expected username to be %s but got %s", username, u.Username)
	}
	if !u.CheckPassword(password) {
		t.Fatalf("Expected password to match %s but it didn't", password)
	}
}
