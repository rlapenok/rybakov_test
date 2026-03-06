package pg

import "time"

// PgConnectionConfig - interface for connection configuration
type PgConnectionConfig struct {
	Host          string
	Port          uint16
	User          string
	Password      string
	Database      string
	Schema        string
	SSLMode       string
	SSLCert       string
	SSLKey        string
	SSLRoot       string
	MigrationPath string
}

// PoolConfig - interface for pool configuration
type PgPoolConfig struct {
	PgConnectionConfig
	MinConns                 int32
	MaxConns                 int32
	MaxConnLifetime          time.Duration
	MaxConnIdleTime          time.Duration
	MaxConnKeepAliveTime     time.Duration
	MaxConnKeepAliveCount    int
	MaxConnKeepAliveInterval time.Duration
}
