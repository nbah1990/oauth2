package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"oauth-proxy/internal/oauth-proxy/config"
	"oauth-proxy/internal/oauth-proxy/infrastructure/database"
	"oauth-proxy/internal/oauth-proxy/migrations"
	"oauth-proxy/internal/server"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
}

func main() {
	cfg := config.New()
	db, err := database.MysqlConnector{}.CreateConnection(cfg.DatabaseConfig)
	if err != nil {
		log.Fatalf(`DB connection fatal error: %s`, err.Error())
	}

	initMigrations(cfg, db)
	initServer(cfg, db)
}

func initMigrations(cfg *config.Config, db *sql.DB) {
	mm := migrations.MigrationManager{DB: db}
	if cfg.ExecuteMigrations {
		log.Println(`Running initMigrations...`)

		err := mm.MigrateFromSource(fmt.Sprintf(`file://%s`, cfg.MigrationsPath))
		if err != nil {
			log.Fatalf(`Migration fatal error: %s`, err.Error())
		}

		log.Println(`Migration completed.`)
	}

	err := mm.MigrateRootClient(cfg.RootClientConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func initServer(cfg *config.Config, db *sql.DB) {
	as := server.APIServer{Config: cfg, Database: db}

	err := as.Init()
	if err != nil {
		log.Fatalf(`API server not started: %s`, err.Error())
	}
}
