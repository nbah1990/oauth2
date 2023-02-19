package entities

import (
	"github.com/docker/distribution/uuid"
	"golang.org/x/exp/slices"
	"oauth-proxy/internal/oauth-proxy/enums"
	"oauth-proxy/internal/oauth-proxy/services/random"
	"time"
)

type Client struct {
	ID        string    `json:"id"`
	Secret    string    `json:"secret"`
	Name      string    `json:"name"`
	Scopes    []string  `json:"scopes"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c Client) HasScope(sc string) bool {
	return slices.IndexFunc(c.Scopes, func(csc string) bool { return csc == sc || csc == string(enums.All) }) != -1
}

func (c Client) CheckSecret(s string) bool {
	return c.Secret == s
}

func CreateClient(name string, scopes []string) *Client {
	return &Client{
		ID:        uuid.Generate().String(),
		Secret:    random.MakeString(80),
		Name:      name,
		Scopes:    scopes,
		Revoked:   false,
		CreatedAt: time.Now().Truncate(time.Millisecond),
		UpdatedAt: time.Now().Truncate(time.Millisecond),
	}
}
