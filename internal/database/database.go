package database

import (
	"fmt"
	"time"

	"github.com/maxx-2345/notes-app-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

type Repository[T any] struct {
	db *Database
}

func NewRepository[T any](db *Database) *Repository[T] {
	return &Repository[T]{db: db}
}

func Connect(cfg *config.DatabaseConfig) (*Database, error) {
	fmt.Printf("Connecting to database: host=%s, port=%s, dbname=%s\n", cfg.Host, cfg.Port, cfg.Name)

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

	fmt.Println("Database connection established and pool tuned")

	return &Database{DB: db}, nil

}


