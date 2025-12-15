package db

import (
	"fmt"
	"log"
	"time"

	"pulse/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DatabaseURL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to database")
	return db, nil
}

// HealthCheck pings the database to verify connectivity
func HealthCheck(db *gorm.DB) error {
	var tmp int
	if err := db.Raw("SELECT 1").Scan(&tmp).Error; err != nil {
		return err
	}
	return nil
}
