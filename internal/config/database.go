package config

import (
	"context"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type DatabaseConfig struct {
	Host         string `env:"DB_HOST" envDefault:"localhost"`
	Port         string `env:"DB_PORT" envDefault:"5432"`
	User         string `env:"DB_USERNAME" envDefault:"postgres"`
	Password     string `env:"DB_PASSWORD,required"`
	Name         string `env:"DB_NAME" envDefault:"postgres"`
	Schema       string `env:"DB_SCHEMA" envDefault:"public"`
	SSLMode      string `env:"SSLMODE" envDefault:"disable"`
	MaxOpenConns int    `env:"MAX_OPEN_CONNS" envDefault:"25"` // Maximum number of open connections to the database
	MaxIdleConns int    `env:"MAX_IDLE_CONNS" envDefault:"5"`  // Maximum number of idle connections in the pool
	// Maximum amount of time a connection may be reused (e.g., "5m", "1h")
	MaxLifetime string `env:"MAX_LIFETIME" envDefault:"5m"`
	// Maximum amount of time a connection may be idle (e.g., "10m", "30m")
	MaxIdleTime string `env:"MAX_IDLE_TIME" envDefault:"10m"`
	// GormLogLevel: silent, error, warn, info (see gorm.io/gorm/logger.LogLevel).
	GormLogLevel string `env:"GORM_LOG_LEVEL" envDefault:"info"`
	DSN          string
}

// LoadDatabaseConfig loads database configuration from environment variables
func LoadDatabaseConfig(ctx context.Context) (*DatabaseConfig, error) {
	cfg := &DatabaseConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	cfg.DSN =
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode, cfg.Schema,
		)
	return cfg, nil
}
