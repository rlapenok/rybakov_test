package pg

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rlapenok/rybakov_test/pkg/db"

	pgxMigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Pool - connection pool to PostgreSQL
type Pool struct {
	pool *pgxpool.Pool
}

// NewPool - create new pool
func NewPool(ctx context.Context, config *PgPoolConfig) (*Pool, error) {

	// Parse config
	poolConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, errPoolCreate(err)
	}

	// Get connection config
	conn := &poolConfig.ConnConfig.Config

	// Set connection config
	conn.Host = config.PgConnectionConfig.Host
	conn.Port = config.PgConnectionConfig.Port
	conn.User = config.PgConnectionConfig.User
	conn.Password = config.PgConnectionConfig.Password
	conn.Database = config.PgConnectionConfig.Database

	// Set search path
	if config.PgConnectionConfig.Schema != "" {
		conn.RuntimeParams["search_path"] = config.PgConnectionConfig.Schema
	}

	// Set TLS config
	if config.PgConnectionConfig.SSLMode != "disable" && config.PgConnectionConfig.SSLCert != "" && config.PgConnectionConfig.SSLKey != "" {
		tlsCfg, err := db.BuildTLSConfig(
			config.PgConnectionConfig.SSLRoot,
			config.PgConnectionConfig.SSLCert,
			config.PgConnectionConfig.SSLKey,
		)
		if err != nil {
			return nil, err
		}
		conn.TLSConfig = tlsCfg
	}

	// Set pool config
	poolConfig.MinConns = config.MinConns
	poolConfig.MaxConns = config.MaxConns
	poolConfig.MaxConnLifetime = config.MaxConnLifetime
	poolConfig.MaxConnIdleTime = config.MaxConnIdleTime

	// Create pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, errPoolCreate(err)
	}

	return &Pool{pool: pool}, nil
}

// Migrate - migrate pool
func (p *Pool) Migrate(ctx context.Context, migrationsPath string) error {
	conn := stdlib.OpenDB(*p.pool.Config().ConnConfig)

	driver, err := pgxMigrate.WithInstance(conn, &pgxMigrate.Config{})
	if err != nil {
		return errPoolMigrate(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"pgx",
		driver,
	)
	if err != nil {
		return errPoolMigrate(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errPoolMigrate(err)
	}

	if err := conn.Close(); err != nil {
		return errPoolMigrate(err)
	}

	return nil
}

// Pgx - get pgx pool
func (p *Pool) Pgx() *pgxpool.Pool {
	return p.pool
}

// Close - close pool
func (p *Pool) Close() {
	p.pool.Close()
}
