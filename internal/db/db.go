package db

import (
	"fmt"
	"os"

	"github.com/sokosplit/sokosplit-core-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect opens a GORM connection using the DATABASE_URL env var, e.g.:
// postgres://user:pass@localhost:5432/sokosplit?sslmode=disable
func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}
	return db, nil
}

// AutoMigrate creates/updates tables for all SokoSplit models.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.SplitList{}, &models.Recipient{})
}
