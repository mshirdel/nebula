package postgres

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
)

func NewPostgres(cfg *Config) (*gorm.DB, error) {
	dsn := cfg.DSN()

	postgres, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(logrus.StandardLogger(), logger.Config{
			SlowThreshold:             cfg.Logger.SlowThreshold,
			LogLevel:                  cfg.Logger.GormLogLevel(),
			Colorful:                  cfg.Logger.Colorful,
			IgnoreRecordNotFoundError: cfg.Logger.IgnoreRecordNotFoundError,
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("can't open DB with: %w", err)
	}

	postgres.Logger = NewGormLogger(cfg.Logger.SlowThreshold)

	err = postgres.Use(prometheus.New(prometheus.Config{DBName: cfg.DBName, StartServer: false}))
	if err != nil {
		return nil, fmt.Errorf("can't use prometheus plugin: %w", err)
	}

	sqlDB, err := postgres.DB()
	if err != nil {
		return nil, fmt.Errorf("can't get DB from gorm: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.PoolSize)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(cfg.ConnectionLifetime)

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("can't ping DB with: %w", err)
	}

	return postgres, nil
}

func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Errorf("get gorm db: %s", err.Error())
	}

	if err := sqlDB.Close(); err != nil {
		logrus.Errorf("close database: %s", err.Error())
	}
}
