package config

import (
	"errors"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rlapenok/rybakov_test/pkg/db/pg"
	pkgErr "github.com/rlapenok/rybakov_test/pkg/errors"
)

var (
	errLoadConfig = pkgErr.NewError(
		pkgErr.CodeInternalError,
		"Failed to load configuration",
		"FAILED_TO_LOAD_CONFIGURATION",
		nil,
	)
)

// Config is a struct that contains the configuration
type Config struct {
	Server   *ServerConfig
	Logger   *LoggerConfig
	Database *DatabaseConfig
	Auth     *AuthConfig
}

// LoggerConfig is a struct that contains the logger configuration
type LoggerConfig struct {
	Level  string `env:"LOGGER_LEVEL"`
	Format string `env:"LOGGER_FORMAT"`
}

// ServerConfig is a struct that contains the server configuration
type ServerConfig struct {
	Port uint16 `env:"HTTP_SERVER_PORT" envDefault:"8080"`
}

// AuthConfig is a struct that contains auth configuration.
type AuthConfig struct {
	BearerToken string `env:"AUTH_BEARER_TOKEN" envDefault:"dev-token"`
}

// DatabaseConfig is a struct that contains database configuration.
type DatabaseConfig struct {
	Host            string        `env:"PG_HOST" envDefault:"localhost"`
	Port            uint16        `env:"PG_PORT" envDefault:"5432"`
	User            string        `env:"PG_USER" envDefault:"postgres"`
	Password        string        `env:"PG_PASSWORD" envDefault:"postgres"`
	Database        string        `env:"PG_DATABASE" envDefault:"postgres"`
	Schema          string        `env:"PG_SCHEMA"`
	SSLMode         string        `env:"PG_SSL_MODE" envDefault:"disable"`
	SSLCert         string        `env:"PG_SSL_CERT"`
	SSLKey          string        `env:"PG_SSL_KEY"`
	SSLRoot         string        `env:"PG_SSL_ROOT"`
	MigrationPath   string        `env:"PG_MIGRATION_PATH" envDefault:"migrations"`
	MinConns        int32         `env:"PG_MIN_CONNS" envDefault:"1"`
	MaxConns        int32         `env:"PG_MAX_CONNS" envDefault:"4"`
	MaxConnLifetime time.Duration `env:"PG_MAX_CONN_LIFETIME" envDefault:"30m"`
	MaxConnIdleTime time.Duration `env:"PG_MAX_CONN_IDLE_TIME" envDefault:"5m"`
}

// LoadConfig loads the configuration from the environment variables
func LoadConfig() (*Config, error) {
	// Initialize the configuration
	var (
		serverConfig   ServerConfig
		loggerConfig   LoggerConfig
		databaseConfig DatabaseConfig
		authConfig     AuthConfig
	)

	// Load the environment variables from the .env file
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, errLoadConfig.WithMessage(err.Error())
	}

	// Parse the server configuration
	if err := env.Parse(&serverConfig); err != nil {
		return nil, errLoadConfig.WithMessage(err.Error())
	}

	// Parse the logger configuration
	if err := env.Parse(&loggerConfig); err != nil {
		return nil, errLoadConfig.WithMessage(err.Error())
	}

	// Parse the database configuration
	if err := env.Parse(&databaseConfig); err != nil {
		return nil, errLoadConfig.WithMessage(err.Error())
	}

	// Parse the auth configuration
	if err := env.Parse(&authConfig); err != nil {
		return nil, errLoadConfig.WithMessage(err.Error())
	}

	// Return the configuration
	return &Config{
		Server: &serverConfig,
		Logger: &loggerConfig,
		Database: &DatabaseConfig{
			Host:            databaseConfig.Host,
			Port:            databaseConfig.Port,
			User:            databaseConfig.User,
			Password:        databaseConfig.Password,
			Database:        databaseConfig.Database,
			Schema:          databaseConfig.Schema,
			SSLMode:         databaseConfig.SSLMode,
			SSLCert:         databaseConfig.SSLCert,
			SSLKey:          databaseConfig.SSLKey,
			SSLRoot:         databaseConfig.SSLRoot,
			MigrationPath:   databaseConfig.MigrationPath,
			MinConns:        databaseConfig.MinConns,
			MaxConns:        databaseConfig.MaxConns,
			MaxConnLifetime: databaseConfig.MaxConnLifetime,
			MaxConnIdleTime: databaseConfig.MaxConnIdleTime,
		},
		Auth: &authConfig,
	}, nil
}

// ToPgPoolConfig converts database config to postgres pool config.
func (c *DatabaseConfig) ToPgPoolConfig() *pg.PgPoolConfig {
	return &pg.PgPoolConfig{
		PgConnectionConfig: pg.PgConnectionConfig{
			Host:          c.Host,
			Port:          c.Port,
			User:          c.User,
			Password:      c.Password,
			Database:      c.Database,
			Schema:        c.Schema,
			SSLMode:       c.SSLMode,
			SSLCert:       c.SSLCert,
			SSLKey:        c.SSLKey,
			SSLRoot:       c.SSLRoot,
			MigrationPath: c.MigrationPath,
		},
		MinConns:        c.MinConns,
		MaxConns:        c.MaxConns,
		MaxConnLifetime: c.MaxConnLifetime,
		MaxConnIdleTime: c.MaxConnIdleTime,
	}
}
