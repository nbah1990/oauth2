package repositories

import (
	"database/sql"
	"oauth-proxy/internal/oauth-proxy/entities"
	dbErrors "oauth-proxy/internal/oauth-proxy/infrastructure/database/errors"
)

type UserRepositoryI interface {
	FindByUsername(username string) (u *entities.User, err error)
	GetByID(id string) (u *entities.User, err error)
	Persist(u *entities.User) error
}

type SQLUserRepository struct {
	DB *sql.DB
}

func (ur SQLUserRepository) FindByUsername(username string) (u *entities.User, err error) {
	u = &entities.User{Username: username}
	err = ur.DB.QueryRow("SELECT u.id, u.password FROM users u WHERE u.username = ?", username).Scan(&u.ID, &u.Password)
	if err != nil {
		if err.Error() == string(dbErrors.NoRows) {
			return nil, dbErrors.NewNotFoundError(`user`)
		}
		return nil, err
	}

	return u, nil
}

func (ur SQLUserRepository) Persist(u *entities.User) error {
	_, err := ur.GetByID(u.ID)

	if err != nil {
		err = ur.insert(u)
	} else {
		err = ur.update(u)
	}

	if err != nil {
		return err
	}

	return nil
}

func (ur SQLUserRepository) GetByID(id string) (u *entities.User, err error) {
	u = &entities.User{ID: id}
	err = ur.DB.QueryRow("SELECT u.id, u.password, u.username FROM users u WHERE u.id = ?", id).Scan(&u.ID, &u.Password, &u.Username)
	if err != nil {
		if err.Error() == string(dbErrors.NoRows) {
			return nil, dbErrors.NewNotFoundError(`user`)
		}
		return nil, err
	}

	return u, nil
}

func (ur SQLUserRepository) insert(u *entities.User) error {
	var eID *string
	if len(u.ExternalID) > 0 {
		eID = &u.ExternalID
	} else {
		eID = nil
	}

	_, err := ur.DB.Exec("INSERT INTO users (id, username, password, external_id) VALUES (?, ?, ?, ?)", u.ID, u.Username, u.Password, eID)
	if err != nil {
		return err
	}
	return nil
}

func (ur SQLUserRepository) update(u *entities.User) error {
	var eID *string
	if len(u.ExternalID) > 0 {
		eID = &u.ExternalID
	} else {
		eID = nil
	}

	_, err := ur.DB.Exec("update users set username = ?, password = ?, external_id = ? where id = ?", u.Username, u.Password, eID, u.ID)
	if err != nil {
		return err
	}
	return nil
}
