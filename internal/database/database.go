package database

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/maxx-2345/notes-app-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.DatabaseConfig, logger *slog.Logger) (*gorm.DB, error) {
	logger.Info("Connecting to database", "host", cfg.Host, "port", cfg.Port, "dbname", cfg.Name)

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}

	maxLifetime, _ := time.ParseDuration(cfg.MaxLifetime)
	maxIdleTime, _ := time.ParseDuration(cfg.MaxIdleTime)

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(maxLifetime)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)

	logger.Info("Database connection established and pool tuned")

	return db, nil

}
