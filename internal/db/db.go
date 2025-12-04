package db

import (
	"fmt"
	"log"

	"pulse/internal/config"
	"pulse/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DatabaseURL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to database")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	log.Println("Running migrations...")
	if err := db.AutoMigrate(
		&models.Project{},
		&models.Tag{},
		&models.Region{},
		&models.Check{},
		&models.CheckRun{},
	); err != nil {
		return err
	}

	return nil
}

func Seed(db *gorm.DB) error {
	if err := seedRegions(db); err != nil {
		return err
	}
	return nil
}

func seedRegions(db *gorm.DB) error {
	regions := []models.Region{
		{Name: "Asia Pacific", Code: "apac"},
	}

	for _, region := range regions {
		var existing models.Region
		err := db.Where("code = ?", region.Code).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&region).Error; err != nil {
				log.Printf("Failed to seed region %s: %v", region.Name, err)
			} else {
				log.Printf("Seeded region %s", region.Name)
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

// HealthCheck pings the database to verify connectivity
func HealthCheck(db *gorm.DB) error {
	var tmp int
	if err := db.Raw("SELECT 1").Scan(&tmp).Error; err != nil {
		return err
	}
	return nil
}
