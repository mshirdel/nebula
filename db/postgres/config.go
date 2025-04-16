package postgres

import (
	"fmt"
	"time"
)

type Config struct {
	Host               string         `mapstructure:"host" validate:"required"`
	Port               int            `mapstructure:"port" validate:"required"`
	Username           string         `mapstructure:"user" validate:"required"`
	Password           string         `mapstructure:"password" validate:"required"`
	DBName             string         `mapstructure:"dbname" validate:"required"`
	ConnectionTimeout  time.Duration  `mapstructure:"connection-timeout" validate:"required"`
	ConnectionLifetime time.Duration  `mapstructure:"connection-lifetime" validate:"required"`
	PoolSize           int            `mapstructure:"pool-size" validate:"required"`
	MaxIdleConnections int            `mapstructure:"max-idle-connections" validate:"required"`
	Logger             DatabaseLogger `mapstructure:"logger" validate:"required"`
}

type DatabaseLogger struct {
	SlowThreshold time.Duration `mapstructure:"slow-threshold" validate:"required"`
}

func (d *Config) DSN() string {
	return d.makeDSN(true)
}

func (d *Config) MigrationDSN() string {
	return d.makeDSN(false)
}

func (d *Config) makeDSN(dbSelected bool) string {
	var dbName string
	if dbSelected {
		dbName = fmt.Sprintf(" database=%s", d.DBName)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s%s connect_timeout=%d sslmode=disable",
		d.Host,
		d.Port,
		d.Username,
		d.Password,
		dbName,
		int(d.ConnectionTimeout.Seconds()),
	)

	return dsn
}
