package migrations

import (
	"pulse/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512102145_add_region_flag",
		Migrate: func(tx *gorm.DB) error {
			return tx.Migrator().AddColumn(&models.Region{}, "flag")
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&models.Region{}, "flag")
		},
	})
}
