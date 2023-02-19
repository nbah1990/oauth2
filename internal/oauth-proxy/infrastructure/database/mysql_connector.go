package database

import (
	"database/sql"
	"fmt"
	"oauth-proxy/internal/oauth-proxy/config"
)

type MysqlConnector struct {
}

func (mc MysqlConnector) CreateConnection(cfg config.DatabaseConfig) (db *sql.DB, err error) {
	cs := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DatabaseName,
	)
	db, err = sql.Open("mysql", cs)
	if err != nil {
		return db, err
	}

	return db, nil
}
