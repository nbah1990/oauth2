package grant

import (
	"oauth-proxy/internal/oauth-proxy/config"
	"oauth-proxy/internal/oauth-proxy/entities"
	"oauth-proxy/internal/oauth-proxy/infrastructure/repositories"
	"time"
)

type Grant interface {
	IssueAccessToken(c *entities.Client, u *entities.User) (at *entities.AccessToken, err error)
	IssueRefreshToken(at *entities.AccessToken) (rt *entities.RefreshToken, err error)

	GetConfig() *config.Config
}

type AbstractGrant struct {
	accessTokenRepository  repositories.AccessTokenRepositoryI
	refreshTokenRepository repositories.RefreshTokenRepositoryI

	config *config.Config
}

func NewAbstractGrant(atr repositories.AccessTokenRepositoryI, rtr repositories.RefreshTokenRepositoryI, config *config.Config) AbstractGrant {
	return AbstractGrant{
		accessTokenRepository:  atr,
		refreshTokenRepository: rtr,
		config:                 config,
	}
}

func (g AbstractGrant) IssueAccessToken(c *entities.Client, u *entities.User) (at *entities.AccessToken, err error) {
	id := g.accessTokenRepository.GenerateNewID()
	at = &entities.AccessToken{
		ID:        id,
		ClientID:  c.ID,
		UserID:    u.ID,
		ExpiresAt: time.Now().Add(g.config.AccessTokenExpirationPeriod).Truncate(time.Millisecond),
	}
	err = g.accessTokenRepository.Persist(at)

	if err != nil {
		return nil, err
	}

	return at, nil
}

func (g AbstractGrant) IssueRefreshToken(at *entities.AccessToken) (rt *entities.RefreshToken, err error) {
	id := g.refreshTokenRepository.GenerateNewID()
	rt = &entities.RefreshToken{
		ID:            id,
		AccessTokenID: at.ID,
		ExpiresAt:     time.Now().Add(g.config.RefreshTokenExpirationPeriod).Truncate(time.Millisecond),
	}
	err = g.refreshTokenRepository.Persist(rt)

	if err != nil {
		return nil, err
	}

	return rt, nil
}

func (g AbstractGrant) GetConfig() *config.Config {
	return g.config
}
