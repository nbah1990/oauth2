package repositories

import (
	"database/sql"
	"oauth-proxy/internal/oauth-proxy/entities"
	dbErrors "oauth-proxy/internal/oauth-proxy/infrastructure/database/errors"
	"strings"
	"time"
)

type ClientRepositoryI interface {
	GetByID(id string) (cli *entities.Client, err error)
	Persist(cli *entities.Client) error
}

type SQLClientRepository struct {
	DB *sql.DB
}

func (cr SQLClientRepository) GetByID(id string) (cli *entities.Client, err error) {
	cli = &entities.Client{ID: id}
	var cAt string
	var uAt string
	var sc string

	err = cr.DB.QueryRow("SELECT oc.secret, oc.name, oc.revoked, oc.scopes, oc.created_at, oc.updated_at FROM oauth_clients oc WHERE id = ?", id).Scan(&cli.Secret, &cli.Name, &cli.Revoked, &sc, &cAt, &uAt)

	if err != nil {
		if err.Error() == string(dbErrors.NoRows) {
			return nil, dbErrors.NewNotFoundError(`client`)
		}
		return nil, err
	}

	cli.Scopes = strings.Fields(sc)
	cli.CreatedAt, _ = time.Parse(time.RFC3339, cAt)
	cli.UpdatedAt, _ = time.Parse(time.RFC3339, uAt)

	return cli, nil
}

func (cr SQLClientRepository) Persist(cli *entities.Client) error {
	_, err := cr.GetByID(cli.ID)
	if err != nil {
		err = cr.insert(cli)
	} else {
		err = cr.update(cli)
	}

	if err != nil {
		return err
	}

	return err
}

func (cr SQLClientRepository) insert(cli *entities.Client) error {
	_, err := cr.DB.Exec("INSERT INTO oauth_clients (id, secret, name, revoked, scopes, created_at, updated_at, redirect) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", cli.ID, cli.Secret, cli.Name, cli.Revoked, strings.Join(cli.Scopes, " "), cli.CreatedAt, cli.UpdatedAt, nil)
	if err != nil {
		return err
	}
	return nil
}

func (cr SQLClientRepository) update(cli *entities.Client) error {
	_, err := cr.DB.Exec("update oauth_clients set secret = ?, name = ?, revoked=?, scopes=?, created_at = ?, updated_at = ? where id = ?", cli.Secret, cli.Name, cli.Revoked, strings.Join(cli.Scopes, " "), cli.CreatedAt, cli.UpdatedAt, cli.ID)
	if err != nil {
		return err
	}
	return nil
}
