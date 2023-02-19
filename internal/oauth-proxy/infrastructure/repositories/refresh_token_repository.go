package repositories

import (
	"database/sql"
	"oauth-proxy/internal/oauth-proxy/entities"
	dbErrors "oauth-proxy/internal/oauth-proxy/infrastructure/database/errors"
	"oauth-proxy/internal/oauth-proxy/services/random"
	"time"
)

type RefreshTokenRepositoryI interface {
	GenerateNewID() string
	Persist(rt *entities.RefreshToken) error
	GetByID(id string) (rt *entities.RefreshToken, err error)
	Revoke(rt *entities.RefreshToken) error
}

type SQLRefreshTokenRepository struct {
	DB *sql.DB
}

func (rtr SQLRefreshTokenRepository) GenerateNewID() string {
	return random.MakeString(100)
}

func (rtr SQLRefreshTokenRepository) Persist(rt *entities.RefreshToken) error {
	_, err := rtr.GetByID(rt.ID)
	if err != nil {
		err = rtr.insert(rt)
	} else {
		err = rtr.update(rt)
	}

	if err != nil {
		return err
	}

	return nil
}

func (rtr SQLRefreshTokenRepository) GetByID(id string) (rt *entities.RefreshToken, err error) {
	rt = &entities.RefreshToken{ID: id}
	var eAt string

	err = rtr.DB.QueryRow("SELECT rt.id, rt.access_token_id, rt.revoked, rt.expires_at FROM oauth_refresh_tokens rt WHERE id = ?", id).Scan(&rt.ID, &rt.AccessTokenID, &rt.Revoked, &eAt)

	if err != nil {
		if err.Error() == string(dbErrors.NoRows) {
			return nil, dbErrors.NewNotFoundError(`refresh_token`)
		}
		return nil, err
	}

	rt.ExpiresAt, _ = time.Parse(time.RFC3339, eAt)

	return rt, nil
}

func (rtr SQLRefreshTokenRepository) Revoke(rt *entities.RefreshToken) error {
	rt.Revoked = true
	return rtr.update(rt)
}

func (rtr SQLRefreshTokenRepository) insert(rt *entities.RefreshToken) error {
	_, err := rtr.DB.Exec("INSERT INTO oauth_refresh_tokens (access_token_id, revoked, expires_at, id) VALUES (?, ?, ?, ?)", rt.AccessTokenID, rt.Revoked, rt.ExpiresAt, rt.ID)
	if err != nil {
		return err
	}
	return nil
}

func (rtr SQLRefreshTokenRepository) update(rt *entities.RefreshToken) error {
	_, err := rtr.DB.Exec("update oauth_refresh_tokens set access_token_id = ?, revoked = ?, expires_at = ? where id = ?", rt.AccessTokenID, rt.Revoked, rt.ExpiresAt, rt.ID)
	if err != nil {
		return err
	}
	return nil
}
