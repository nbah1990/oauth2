package entities

import (
	"github.com/docker/distribution/uuid"
	"oauth-proxy/internal/oauth-proxy/services/hash"
)

type User struct {
	ID         string
	Username   string
	Password   string
	ExternalID string
}

func (u User) CheckPassword(pass string) bool {
	return hash.IsHashSame(pass, u.Password)
}

func CreateUser(username string, password string, externalID string) *User {
	pH, _ := hash.BcryptHash(password)
	u := &User{
		ID:       uuid.Generate().String(),
		Username: username,
		Password: pH,
	}

	if len(externalID) > 0 {
		u.ExternalID = externalID
	}

	return u
}
