package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"log"
	"oauth-proxy/internal/oauth-proxy/config"
	"oauth-proxy/internal/oauth-proxy/entities"
	"oauth-proxy/internal/oauth-proxy/enums"
	"oauth-proxy/internal/oauth-proxy/infrastructure/repositories"
)

type MigrationManager struct {
	DB *sql.DB
}

func (mm MigrationManager) MigrateFromSource(source string) error {
	driver, err := mysql.WithInstance(mm.DB, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		source,
		"",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println(`Nothing to migrate.`)
		} else {
			return err
		}
	}

	return nil
}

func (mm MigrationManager) MigrateRootClient(rcc config.RootClientConfig) error {
	scr := repositories.SQLClientRepository{DB: mm.DB}
	client, err := scr.GetByID(rcc.ClientID)

	if rcc.Enabled {
		if err != nil {
			client = entities.CreateClient(`root_client`, []string{})
		}

		client.ID = rcc.ClientID
		client.Secret = rcc.Secret
		client.Scopes = []string{string(enums.All)}
		client.Revoked = false

		err = scr.Persist(client)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf(`can not create a root client:%s`, err.Error()))
		}
	} else if err == nil {
		client.Revoked = true
		err = scr.Persist(client)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf(`can not revoke a root client:%s`, err.Error()))
		}
	}
	return nil
}
