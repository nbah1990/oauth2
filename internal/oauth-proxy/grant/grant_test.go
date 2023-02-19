package grant

import (
	"errors"
	"github.com/docker/distribution/uuid"
	"oauth-proxy/internal/oauth-proxy/config"
	"oauth-proxy/internal/oauth-proxy/entities"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAccessTokenRepository struct {
	mock.Mock
}

func (m *mockAccessTokenRepository) GetByID(id string) (client *entities.AccessToken, err error) {
	args := m.Called(id)
	return args.Get(0).(*entities.AccessToken), args.Error(1)
}

func (m *mockAccessTokenRepository) Revoke(at *entities.AccessToken) error {
	args := m.Called(at)
	return args.Error(0)
}

func (m *mockAccessTokenRepository) GenerateNewID() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockAccessTokenRepository) Persist(at *entities.AccessToken) error {
	args := m.Called(at)
	return args.Error(0)
}

func TestGrantIssueAccessToken_Success(t *testing.T) {
	mockAccessTokenRepository := new(mockAccessTokenRepository)
	cfg := &config.Config{
		AccessTokenExpirationPeriod: 10 * time.Minute,
	}
	grant := NewAbstractGrant(mockAccessTokenRepository, nil, cfg)

	client := &entities.Client{
		ID: "client_id",
	}

	user := &entities.User{
		ID: "user_id",
	}

	expectedAccessToken := &entities.AccessToken{
		ID:        "access_token_id",
		ClientID:  client.ID,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(cfg.AccessTokenExpirationPeriod).Truncate(time.Millisecond),
	}

	mockAccessTokenRepository.On("GenerateNewID").Return("access_token_id")
	mockAccessTokenRepository.On("Persist", expectedAccessToken).Return(nil)

	accessToken, err := grant.IssueAccessToken(client, user)

	assert.Nil(t, err)
	assert.Equal(t, expectedAccessToken, accessToken)
	mockAccessTokenRepository.AssertExpectations(t)
}

func TestGrantIssueAccessToken_ErrorFromRepository(t *testing.T) {
	cfg := &config.Config{
		AccessTokenExpirationPeriod: time.Hour,
	}

	accessTokenRepositoryMock := new(mockAccessTokenRepository)
	accessTokenRepositoryMock.On("GenerateNewID").Return("1")
	accessTokenRepositoryMock.On("Persist", mock.Anything).Return(errors.New("error while persisting access token"))

	grant := NewAbstractGrant(accessTokenRepositoryMock, nil, cfg)
	client := &entities.Client{ID: "1"}
	user := &entities.User{ID: "1"}

	accessToken, err := grant.IssueAccessToken(client, user)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
	if accessToken != nil {
		t.Errorf("Expected access token to be nil, but got %v", accessToken)
	}

	accessTokenRepositoryMock.AssertExpectations(t)
}

type mockRefreshTokenRepository struct {
	mock.Mock
}

func (m *mockRefreshTokenRepository) GetByID(id string) (client *entities.RefreshToken, err error) {
	args := m.Called(id)
	return args.Get(0).(*entities.RefreshToken), args.Error(1)
}

func (m *mockRefreshTokenRepository) Revoke(rt *entities.RefreshToken) error {
	args := m.Called(rt)
	return args.Error(0)
}

func (m *mockRefreshTokenRepository) GenerateNewID() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockRefreshTokenRepository) Persist(at *entities.RefreshToken) error {
	args := m.Called(at)
	return args.Error(0)
}

func TestGrantIssueRefreshToken_Success(t *testing.T) {
	mockRefreshTokenRepository := new(mockRefreshTokenRepository)
	cfg := &config.Config{
		RefreshTokenExpirationPeriod: 1 * time.Hour,
	}
	grant := NewAbstractGrant(nil, mockRefreshTokenRepository, cfg)

	at := &entities.AccessToken{
		ID: uuid.Generate().String(),
	}

	expectedRefreshToken := &entities.RefreshToken{
		ID:            "access_token_id",
		AccessTokenID: at.ID,
		ExpiresAt:     time.Now().Add(cfg.RefreshTokenExpirationPeriod).Truncate(time.Millisecond),
	}

	mockRefreshTokenRepository.On("GenerateNewID").Return("access_token_id")
	mockRefreshTokenRepository.On("Persist", expectedRefreshToken).Return(nil)

	refreshToken, err := grant.IssueRefreshToken(at)

	assert.Nil(t, err)
	assert.Equal(t, expectedRefreshToken, refreshToken)
	mockRefreshTokenRepository.AssertExpectations(t)
}

func TestGrantIssueRefreshToken_ErrorFromRepository(t *testing.T) {
	cfg := &config.Config{
		RefreshTokenExpirationPeriod: time.Hour,
	}

	refreshTokenRepositoryMock := new(mockRefreshTokenRepository)
	refreshTokenRepositoryMock.On("GenerateNewID").Return("1")
	refreshTokenRepositoryMock.On("Persist", mock.Anything).Return(errors.New("error while persisting access token"))

	grant := NewAbstractGrant(nil, refreshTokenRepositoryMock, cfg)
	at := &entities.AccessToken{
		ID: uuid.Generate().String(),
	}

	refreshToken, err := grant.IssueRefreshToken(at)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
	if refreshToken != nil {
		t.Errorf("Expected access token to be nil, but got %v", refreshToken)
	}

	refreshTokenRepositoryMock.AssertExpectations(t)
}

func TestGrantGetConfig(t *testing.T) {
	expectedConfig := &config.Config{
		AccessTokenExpirationPeriod:  time.Minute * 30,
		RefreshTokenExpirationPeriod: time.Hour * 24,
	}

	abstractGrant := AbstractGrant{config: expectedConfig}

	actualConfig := abstractGrant.GetConfig()

	if !reflect.DeepEqual(expectedConfig, actualConfig) {
		t.Errorf("Expected config %v, but got %v", expectedConfig, actualConfig)
	}
}
