package config

import (
	"os"
	"time"
)

const TrueString = `true`

type Config struct {
	ApplicationAddress           string
	AccessTokenExpirationPeriod  time.Duration
	RefreshTokenExpirationPeriod time.Duration

	DatabaseConfig    DatabaseConfig
	RootClientConfig  RootClientConfig
	ExecuteMigrations bool
	MigrationsPath    string
}

func New() *Config {
	atp := GetEnv(`ACCESS_TOKEN_EXPIRATION_PERIOD`, `300s`)
	rtp := GetEnv(`REFRESH_TOKEN_EXPIRATION_PERIOD`, `900s`)

	atpD, err := time.ParseDuration(atp)
	if err != nil {
		panic(err)
	}

	rtpD, err := time.ParseDuration(rtp)
	if err != nil {
		panic(err)
	}

	return &Config{
		ApplicationAddress:           GetEnv(`APP_ADDRESS`, `0.0.0.0:8096`),
		AccessTokenExpirationPeriod:  atpD,
		RefreshTokenExpirationPeriod: rtpD,
		DatabaseConfig:               NewDatabaseConfig(),
		RootClientConfig:             NewRootClientConfig(),
		ExecuteMigrations:            GetEnv(`EXECUTE_MIGRATIONS`, ``) == TrueString,
		MigrationsPath:               GetEnv(`MIGRATIONS_PATH`, ``),
	}
}

func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
