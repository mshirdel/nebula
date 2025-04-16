package postgres

import "time"

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
