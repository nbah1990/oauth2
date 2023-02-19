package grant

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"oauth-proxy/internal/oauth-proxy/config"
	"oauth-proxy/internal/oauth-proxy/entities"
	"oauth-proxy/internal/oauth-proxy/requests"
	"oauth-proxy/internal/oauth-proxy/responses"
	"oauth-proxy/internal/oauth-proxy/services/hash"
	"oauth-proxy/internal/oauth-proxy/services/random"
	"testing"
	"time"
)

type mockClientRepository struct {
	mock.Mock
}

func (m *mockClientRepository) GetByID(id string) (*entities.Client, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (m *mockClientRepository) Persist(client *entities.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) FindByUsername(username string) (*entities.User, error) {
	args := m.Called(username)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *mockUserRepository) GetByID(id string) (*entities.User, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *mockUserRepository) Persist(u *entities.User) error {
	args := m.Called(u)
	return args.Error(0)
}

type mockGrant struct {
	mock.Mock
}

func (m *mockGrant) GetConfig() *config.Config {
	args := m.Called()
	return args.Get(0).(*config.Config)
}

func (m *mockGrant) IssueAccessToken(c *entities.Client, u *entities.User) (at *entities.AccessToken, err error) {
	args := m.Called(c, u)
	return args.Get(0).(*entities.AccessToken), args.Error(1)
}

func (m *mockGrant) IssueRefreshToken(at *entities.AccessToken) (rt *entities.RefreshToken, err error) {
	args := m.Called(at)
	return args.Get(0).(*entities.RefreshToken), args.Error(1)
}

func TestAccessTokenRequest_ClientNotFound(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}

	// Arrange
	pg := PasswordGrant{
		ClientRepository: mcr,
		UserRepository:   mur,
	}
	req := requests.TokenRequest{
		ClientID:  `client-id`,
		Secret:    `secret`,
		Username:  `username`,
		Password:  `password`,
		GrantType: `password`,
	}

	mcr.On(`GetByID`, `client-id`).Return(&entities.Client{}, errors.New(`client not found`))

	// Act
	res, err := pg.AccessTokenRequest(req)

	// Assert
	assert.Nil(t, res)
	assert.Equal(t, `client not found`, err.Error())
}

func TestAccessTokenRequest_ClientNotFound_IncorrectSecret(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}

	// Arrange
	pg := PasswordGrant{
		ClientRepository: mcr,
		UserRepository:   mur,
	}
	req := requests.TokenRequest{
		ClientID:  `client-id`,
		Secret:    `secret`,
		Username:  `username`,
		Password:  `password`,
		GrantType: `password`,
	}

	mcr.On(`GetByID`, `client-id`).Return(&entities.Client{ID: `client-id`, Secret: `secret_2`}, nil)

	// Act
	res, err := pg.AccessTokenRequest(req)

	// Assert
	assert.Nil(t, res)
	assert.Equal(t, `client not found`, err.Error())
}

func TestAccessTokenRequest_UserNotFound(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}

	// Arrange
	pg := PasswordGrant{
		ClientRepository: mcr,
		UserRepository:   mur,
	}
	req := requests.TokenRequest{
		ClientID:  `client-id`,
		Secret:    `secret`,
		Username:  `username`,
		Password:  `password`,
		GrantType: `password`,
	}

	mcr.On(`GetByID`, `client-id`).Return(&entities.Client{ID: `client-id`, Secret: `secret`}, nil)
	mur.On(`FindByUsername`, `username`).Return(&entities.User{}, errors.New(`user not found`))

	// Act
	res, err := pg.AccessTokenRequest(req)

	// Assert
	assert.Nil(t, res)
	assert.Equal(t, `user not found`, err.Error())
}

func TestAccessTokenRequest_UserNotFound_InvalidPassword(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}

	// Arrange
	pg := PasswordGrant{
		ClientRepository: mcr,
		UserRepository:   mur,
	}
	req := requests.TokenRequest{
		ClientID:  `client-id`,
		Secret:    `secret`,
		Username:  `username`,
		Password:  `password`,
		GrantType: `password`,
	}

	mcr.On(`GetByID`, `client-id`).Return(&entities.Client{ID: `client-id`, Secret: `secret`}, nil)

	pass, _ := hash.BcryptHash(`password_2`)
	mur.On(`FindByUsername`, `username`).Return(&entities.User{Username: `username`, Password: pass}, nil)

	// Act
	res, err := pg.AccessTokenRequest(req)

	// Assert
	assert.Nil(t, res)
	assert.Equal(t, `user not found`, err.Error())
}

func TestAccessTokenRequest_InternalError_CanNotIssueAccessToken(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}
	gr := &mockGrant{}

	// Arrange
	pg := PasswordGrant{
		ClientRepository: mcr,
		UserRepository:   mur,
		Grant:            gr,
	}
	req := requests.TokenRequest{
		ClientID:  `client-id`,
		Secret:    `secret`,
		Username:  `username`,
		Password:  `password`,
		GrantType: `password`,
	}

	cli := &entities.Client{ID: `client-id`, Secret: `secret`}
	mcr.On(`GetByID`, `client-id`).Return(cli, nil)

	pass, _ := hash.BcryptHash(`password`)
	u := &entities.User{Username: `username`, Password: pass}
	mur.On(`FindByUsername`, `username`).Return(u, nil)

	gr.On(`IssueAccessToken`, cli, u).Return(&entities.AccessToken{}, errors.New(`internal error`))

	// Act
	res, err := pg.AccessTokenRequest(req)

	// Assert
	assert.Nil(t, res)
	assert.Equal(t, `at internal error`, err.Error())
}

func TestAccessTokenRequest_InternalError_CanNotIssueRefreshToken(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}
	gr := &mockGrant{}

	// Arrange
	pg := PasswordGrant{
		ClientRepository: mcr,
		UserRepository:   mur,
		Grant:            gr,
	}
	req := requests.TokenRequest{
		ClientID:  `client-id`,
		Secret:    `secret`,
		Username:  `username`,
		Password:  `password`,
		GrantType: `password`,
	}

	cli := &entities.Client{ID: `client-id`, Secret: `secret`}
	mcr.On(`GetByID`, `client-id`).Return(cli, nil)

	pass, _ := hash.BcryptHash(`password`)
	u := &entities.User{Username: `username`, Password: pass}
	mur.On(`FindByUsername`, `username`).Return(u, nil)

	at := &entities.AccessToken{ID: random.MakeString(50)}
	gr.On(`IssueAccessToken`, cli, u).Return(at, nil)
	gr.On(`IssueRefreshToken`, at).Return(&entities.RefreshToken{}, errors.New(`internal error`))

	// Act
	res, err := pg.AccessTokenRequest(req)

	// Assert
	assert.Nil(t, res)
	assert.Equal(t, `rt internal error`, err.Error())
}

func TestAccessTokenRequest_Success(t *testing.T) {
	mcr := &mockClientRepository{}
	mur := &mockUserRepository{}
	gr := &mockGrant{}

	// Arrange
	pg := PasswordGrant{
		ClientRepository: mcr,
		UserRepository:   mur,
		Grant:            gr,
	}
	req := requests.TokenRequest{
		ClientID:  `client-id`,
		Secret:    `secret`,
		Username:  `username`,
		Password:  `password`,
		GrantType: `password`,
	}

	cli := &entities.Client{ID: `client-id`, Secret: `secret`}
	mcr.On(`GetByID`, `client-id`).Return(cli, nil)

	pass, _ := hash.BcryptHash(`password`)
	u := &entities.User{Username: `username`, Password: pass}
	mur.On(`FindByUsername`, `username`).Return(u, nil)

	at := &entities.AccessToken{ID: random.MakeString(50)}
	gr.On(`IssueAccessToken`, cli, u).Return(at, nil)

	rt := &entities.RefreshToken{ID: random.MakeString(50)}
	gr.On(`IssueRefreshToken`, at).Return(rt, nil)

	atD, _ := time.ParseDuration(`10s`)
	cfg := &config.Config{AccessTokenExpirationPeriod: atD}
	gr.On(`GetConfig`).Return(cfg)

	exRes := responses.HTTPTokenResponse{
		AccessToken:  at.ID,
		RefreshToken: rt.ID,
		TokenType:    `Bearer`,
		ExpiresIn:    int(cfg.AccessTokenExpirationPeriod.Seconds()),
	}
	// Act
	res, err := pg.AccessTokenRequest(req)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, exRes, res)
}
