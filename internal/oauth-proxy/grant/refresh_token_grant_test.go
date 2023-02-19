package grant

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"oauth-proxy/internal/oauth-proxy/config"
	"oauth-proxy/internal/oauth-proxy/entities"
	"oauth-proxy/internal/oauth-proxy/requests"
	"oauth-proxy/internal/oauth-proxy/responses"
	"oauth-proxy/internal/oauth-proxy/services/random"
	"testing"
	"time"
)

func TestRefreshTokenRequest_ClientNotFound(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}

	// Arrange
	rtg := RefreshTokenGrant{
		ClientRepository: mcr,
		UserRepository:   mur,
	}
	req := requests.RefreshTokenRequest{
		ClientID:  `client-id`,
		Secret:    `secret`,
		GrantType: `refresh_token`,
	}

	mcr.On(`GetByID`, `client-id`).Return(&entities.Client{}, errors.New(`any error`))

	// Act
	res, err := rtg.RefreshTokenRequest(req)

	// Assert
	assert.Nil(t, res)
	assert.Equal(t, `client not found`, err.Error())
}

func TestRefreshTokenRequest_ClientNotFound_InvalidSecret(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}

	// Arrange
	rtg := RefreshTokenGrant{
		ClientRepository: mcr,
		UserRepository:   mur,
	}
	req := requests.RefreshTokenRequest{
		ClientID:  `client-id`,
		Secret:    `secret`,
		GrantType: `refresh_token`,
	}

	cli := &entities.Client{Secret: `secret_2`}
	mcr.On(`GetByID`, `client-id`).Return(cli, nil)

	// Act
	res, err := rtg.RefreshTokenRequest(req)

	// Assert
	assert.Nil(t, res)
	assert.Equal(t, `client not found`, err.Error())
}

func TestRefreshTokenRequest_RefreshTokenNotFound(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}
	rtr := &mockRefreshTokenRepository{}

	// Arrange
	rtg := RefreshTokenGrant{
		ClientRepository:       mcr,
		UserRepository:         mur,
		RefreshTokenRepository: rtr,
	}
	req := requests.RefreshTokenRequest{
		ClientID:     `client-id`,
		Secret:       `secret`,
		RefreshToken: random.MakeString(10),
		GrantType:    `refresh_token`,
	}

	mcr.On(`GetByID`, `client-id`).Return(&entities.Client{ID: `client-id`, Secret: `secret`}, nil)
	rtr.On(`GetByID`, req.RefreshToken).Return(&entities.RefreshToken{}, errors.New(`any error`))

	// Act
	res, err := rtg.RefreshTokenRequest(req)

	// Assert
	assert.Nil(t, res)
	assert.Equal(t, `refresh token not found`, err.Error())
}

func TestRefreshTokenRequest_RefreshTokenExpired(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}
	rtr := &mockRefreshTokenRepository{}

	// Arrange
	rtg := RefreshTokenGrant{
		ClientRepository:       mcr,
		UserRepository:         mur,
		RefreshTokenRepository: rtr,
	}
	req := requests.RefreshTokenRequest{
		ClientID:     `client-id`,
		Secret:       `secret`,
		RefreshToken: random.MakeString(10),
		GrantType:    `refresh_token`,
	}

	mcr.On(`GetByID`, `client-id`).Return(&entities.Client{ID: `client-id`, Secret: `secret`}, nil)
	rtr.On(`GetByID`, req.RefreshToken).Return(&entities.RefreshToken{ExpiresAt: time.Now().Add(-24 * time.Hour)}, nil)

	// Act
	res, err := rtg.RefreshTokenRequest(req)

	// Assert
	assert.Nil(t, res)
	assert.Equal(t, `refresh token expired`, err.Error())
}

func TestRefreshTokenRequest_AccessTokenNotFound(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}
	rtr := &mockRefreshTokenRepository{}
	atr := &mockAccessTokenRepository{}

	// Arrange
	rtg := RefreshTokenGrant{
		ClientRepository:       mcr,
		UserRepository:         mur,
		RefreshTokenRepository: rtr,
		AccessTokenRepository:  atr,
	}
	req := requests.RefreshTokenRequest{
		ClientID:     `client-id`,
		Secret:       `secret`,
		RefreshToken: random.MakeString(10),
		GrantType:    `refresh_token`,
	}

	mcr.On(`GetByID`, `client-id`).Return(&entities.Client{ID: `client-id`, Secret: `secret`}, nil)

	rt := &entities.RefreshToken{ExpiresAt: time.Now().Add(24 * time.Hour), ID: random.MakeString(10), AccessTokenID: random.MakeString(10)}
	rtr.On(`GetByID`, req.RefreshToken).Return(rt, nil)

	atr.On(`GetByID`, rt.AccessTokenID).Return(&entities.AccessToken{}, errors.New(`any error`))

	// Act
	res, err := rtg.RefreshTokenRequest(req)

	// Assert
	assert.Nil(t, res)
	assert.Equal(t, `access token not found`, err.Error())
}

func TestRefreshTokenRequest_UserNotFound(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}
	rtr := &mockRefreshTokenRepository{}
	atr := &mockAccessTokenRepository{}
	gr := &mockGrant{}

	// Arrange
	rtg := RefreshTokenGrant{
		Grant: gr,

		ClientRepository:       mcr,
		UserRepository:         mur,
		RefreshTokenRepository: rtr,
		AccessTokenRepository:  atr,
	}
	req := requests.RefreshTokenRequest{
		ClientID:     `client-id`,
		Secret:       `secret`,
		RefreshToken: random.MakeString(10),
		GrantType:    `refresh_token`,
	}

	mcr.On(`GetByID`, `client-id`).Return(&entities.Client{ID: `client-id`, Secret: `secret`}, nil)

	oldRt := &entities.RefreshToken{ExpiresAt: time.Now().Add(24 * time.Hour), ID: random.MakeString(10), AccessTokenID: random.MakeString(10)}
	rtr.On(`GetByID`, req.RefreshToken).Return(oldRt, nil)

	oldAt := &entities.AccessToken{ID: random.MakeString(10), UserID: `user-id`}
	atr.On(`GetByID`, oldRt.AccessTokenID).Return(oldAt, nil)

	mur.On(`GetByID`, `user-id`).Return(&entities.User{}, errors.New(`any error`))

	// Act
	res, err := rtg.RefreshTokenRequest(req)

	// Assert
	assert.Nil(t, res)
	assert.Equal(t, `user not found`, err.Error())
}

func TestRefreshTokenRequest_InternalError_CanNotIssueAccessToken(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}
	rtr := &mockRefreshTokenRepository{}
	atr := &mockAccessTokenRepository{}
	gr := &mockGrant{}

	// Arrange
	rtg := RefreshTokenGrant{
		Grant: gr,

		ClientRepository:       mcr,
		UserRepository:         mur,
		RefreshTokenRepository: rtr,
		AccessTokenRepository:  atr,
	}
	req := requests.RefreshTokenRequest{
		ClientID:     `client-id`,
		Secret:       `secret`,
		RefreshToken: random.MakeString(10),
		GrantType:    `refresh_token`,
	}

	cli := &entities.Client{ID: `client-id`, Secret: `secret`}
	mcr.On(`GetByID`, `client-id`).Return(cli, nil)

	oldRt := &entities.RefreshToken{ExpiresAt: time.Now().Add(24 * time.Hour), ID: random.MakeString(10), AccessTokenID: random.MakeString(10)}
	rtr.On(`GetByID`, req.RefreshToken).Return(oldRt, nil)

	oldAt := &entities.AccessToken{ID: random.MakeString(10), UserID: `user-id`}
	atr.On(`GetByID`, oldRt.AccessTokenID).Return(oldAt, nil)

	u := &entities.User{ID: `user-id`}
	mur.On(`GetByID`, `user-id`).Return(u, nil)

	gr.On(`IssueAccessToken`, cli, u).Return(&entities.AccessToken{}, errors.New(`any error`))
	// Act
	res, err := rtg.RefreshTokenRequest(req)

	// Assert
	assert.Nil(t, res)
	assert.Equal(t, `at internal error`, err.Error())
}

func TestRefreshTokenRequest_InternalError_CanNotIssueRefreshToken(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}
	rtr := &mockRefreshTokenRepository{}
	atr := &mockAccessTokenRepository{}
	gr := &mockGrant{}

	// Arrange
	rtg := RefreshTokenGrant{
		Grant: gr,

		ClientRepository:       mcr,
		UserRepository:         mur,
		RefreshTokenRepository: rtr,
		AccessTokenRepository:  atr,
	}
	req := requests.RefreshTokenRequest{
		ClientID:     `client-id`,
		Secret:       `secret`,
		RefreshToken: random.MakeString(10),
		GrantType:    `refresh_token`,
	}

	cli := &entities.Client{ID: `client-id`, Secret: `secret`}
	mcr.On(`GetByID`, `client-id`).Return(cli, nil)

	oldRt := &entities.RefreshToken{ExpiresAt: time.Now().Add(24 * time.Hour), ID: random.MakeString(10), AccessTokenID: random.MakeString(10)}
	rtr.On(`GetByID`, req.RefreshToken).Return(oldRt, nil)

	oldAt := &entities.AccessToken{ID: random.MakeString(10), UserID: `user-id`}
	atr.On(`GetByID`, oldRt.AccessTokenID).Return(oldAt, nil)

	u := &entities.User{ID: `user-id`}
	mur.On(`GetByID`, `user-id`).Return(u, nil)

	at := &entities.AccessToken{ID: random.MakeString(50)}
	gr.On(`IssueAccessToken`, cli, u).Return(at, nil)

	gr.On(`IssueRefreshToken`, at).Return(&entities.RefreshToken{}, errors.New(`any error`))

	// Act
	res, err := rtg.RefreshTokenRequest(req)

	// Assert
	assert.Nil(t, res)
	assert.Equal(t, `rt internal error`, err.Error())
}

func TestRefreshTokenRequest_Success(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}
	rtr := &mockRefreshTokenRepository{}
	atr := &mockAccessTokenRepository{}
	gr := &mockGrant{}

	// Arrange
	rtg := RefreshTokenGrant{
		Grant: gr,

		ClientRepository:       mcr,
		UserRepository:         mur,
		RefreshTokenRepository: rtr,
		AccessTokenRepository:  atr,
	}
	req := requests.RefreshTokenRequest{
		ClientID:     `client-id`,
		Secret:       `secret`,
		RefreshToken: random.MakeString(10),
		GrantType:    `refresh_token`,
	}

	cli := &entities.Client{ID: `client-id`, Secret: `secret`}
	mcr.On(`GetByID`, `client-id`).Return(cli, nil)

	oldRt := &entities.RefreshToken{ExpiresAt: time.Now().Add(24 * time.Hour), ID: random.MakeString(10), AccessTokenID: random.MakeString(10)}
	rtr.On(`GetByID`, req.RefreshToken).Return(oldRt, nil)

	oldAt := &entities.AccessToken{ID: random.MakeString(10), UserID: `user-id`}
	atr.On(`GetByID`, oldRt.AccessTokenID).Return(oldAt, nil)

	u := &entities.User{ID: `user-id`}
	mur.On(`GetByID`, `user-id`).Return(u, nil)

	at := &entities.AccessToken{ID: random.MakeString(50)}
	gr.On(`IssueAccessToken`, cli, u).Return(at, nil)

	rt := &entities.RefreshToken{ID: random.MakeString(50)}
	gr.On(`IssueRefreshToken`, at).Return(rt, nil)

	atD, _ := time.ParseDuration(`10s`)
	cfg := &config.Config{AccessTokenExpirationPeriod: atD}
	gr.On(`GetConfig`).Return(cfg)

	atr.On(`Revoke`, oldAt).Return(errors.New(`any error`))
	rtr.On(`Revoke`, oldRt).Return(errors.New(`any error`))

	exRes := responses.HTTPTokenResponse{
		AccessToken:  at.ID,
		RefreshToken: rt.ID,
		TokenType:    `Bearer`,
		ExpiresIn:    int(cfg.AccessTokenExpirationPeriod.Seconds()),
	}
	// Act
	res, err := rtg.RefreshTokenRequest(req)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, exRes, res)
}
