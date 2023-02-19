package repositories

import (
	"database/sql"
	"oauth-proxy/internal/oauth-proxy/entities"
	dbErrors "oauth-proxy/internal/oauth-proxy/infrastructure/database/errors"
	"oauth-proxy/internal/oauth-proxy/services/random"
	"time"
)

type AccessTokenRepositoryI interface {
	GenerateNewID() string
	Persist(at *entities.AccessToken) error
	GetByID(id string) (client *entities.AccessToken, err error)
	Revoke(at *entities.AccessToken) error
}

type SQLAccessTokenRepository struct {
	DB *sql.DB
}

func (atr SQLAccessTokenRepository) GenerateNewID() string {
	return random.MakeString(100)
}

func (atr SQLAccessTokenRepository) Persist(at *entities.AccessToken) error {
	_, err := atr.GetByID(at.ID)
	if err != nil {
		err = atr.insert(at)
	} else {
		err = atr.update(at)
	}

	if err != nil {
		return err
	}

	return nil
}

func (atr SQLAccessTokenRepository) GetByID(id string) (at *entities.AccessToken, err error) {
	at = &entities.AccessToken{ID: id}
	var eAt string

	err = atr.DB.QueryRow("SELECT at.id, at.client_id, at.user_id, at.revoked, at.expires_at FROM oauth_access_tokens at WHERE id = ?", id).Scan(&at.ID, &at.ClientID, &at.UserID, &at.Revoked, &eAt)

	if err != nil {
		if err.Error() == string(dbErrors.NoRows) {
			return nil, dbErrors.NewNotFoundError(`access_token`)
		}
		return nil, err
	}

	at.ExpiresAt, _ = time.Parse(time.RFC3339, eAt)

	return at, nil
}

func (atr SQLAccessTokenRepository) Revoke(at *entities.AccessToken) error {
	at.Revoked = true
	return atr.update(at)
}

func (atr SQLAccessTokenRepository) insert(at *entities.AccessToken) error {
	_, err := atr.DB.Exec("INSERT INTO oauth_access_tokens (id, user_id, client_id, revoked, expires_at) VALUES (?, ?, ?, ?, ?)", at.ID, at.UserID, at.ClientID, at.Revoked, at.ExpiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (atr SQLAccessTokenRepository) update(at *entities.AccessToken) error {
	_, err := atr.DB.Exec("update oauth_access_tokens set user_id = ?, client_id = ?, revoked = ?, expires_at = ? where id = ?", at.UserID, at.ClientID, at.Revoked, at.ExpiresAt, at.ID)
	if err != nil {
		return err
	}
	return nil
}
