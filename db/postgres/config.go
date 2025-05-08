package postgres

import (
	"fmt"
	"time"

	"gorm.io/gorm/logger"
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
	SSLMode            string         `mapstructure:"sslmode" validate:"required"`
	TimeZone           string         `mapstructure:"time-zone" validate:"required"`
	Logger             DatabaseLogger `mapstructure:"logger" validate:"required"`
}

type DatabaseLogger struct {
	SlowThreshold             time.Duration `mapstructure:"slow-threshold" validate:"required"`
	Level                     string        `mapstructure:"level" validate:"required,lowercase,oneof=silent error warn warning info"`
	Colorful                  bool          `mapstructure:"colorful"`
	IgnoreRecordNotFoundError bool          `mapstructure:"ignore-record-not-found-error"`
}

func (l DatabaseLogger) GormLogLevel() logger.LogLevel {
	switch l.Level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn", "warning":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}

func (d *Config) DSN() string {
	return d.makeDSN(true)
}

func (d *Config) MigrationDSN() string {
	return d.makeDSN(false)
}

func (d *Config) makeDSN(_ bool) string {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s connect_timeout=%d",
		d.Host,
		d.Username,
		d.Password,
		d.DBName,
		d.Port,
		d.SSLMode,
		d.TimeZone,
		int(d.ConnectionTimeout.Seconds()),
	)

	return dsn
}
