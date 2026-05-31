package config

import (
	"context"
	"github.com/caarlos0/env/v11"
)

// Config holds all the environment variables for the app.
// We use tags to map the variable names, just like Laravel.
type Config struct {
	Database *DatabaseConfig
}

// Load reads the environment variables and returns a Config struct.
func LoadConfig(ctx context.Context) (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	dbConfig, err := LoadDatabaseConfig(ctx)
	if err != nil {
		return nil, err
	}

	cfg.Database = dbConfig
	
	return cfg, nil
}
