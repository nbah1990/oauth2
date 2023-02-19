package config

type DatabaseConfig struct {
	Host         string
	Port         string
	DatabaseName string
	User         string
	Password     string
}

func NewDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:         GetEnv(`DATABASE_HOST`, ``),
		Port:         GetEnv(`DATABASE_PORT`, ``),
		DatabaseName: GetEnv(`DATABASE_NAME`, ``),
		User:         GetEnv(`DATABASE_USER`, ``),
		Password:     GetEnv(`DATABASE_PASSWORD`, ``),
	}
}
