package migrations

import (
	"pulse/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512060000_add_sessions",
		Migrate: func(tx *gorm.DB) error {
			return tx.Migrator().CreateTable(&models.Session{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&models.Session{})
		},
	})
}
