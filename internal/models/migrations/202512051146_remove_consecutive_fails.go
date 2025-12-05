package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"pulse/internal/models"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512051146_remove_consecutive_fails",
		Migrate: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&models.Check{}, "consecutive_fails")
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().AddColumn(&models.Check{}, "consecutive_fails")
		},
	})
}
