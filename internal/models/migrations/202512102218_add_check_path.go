package migrations

import (
	"pulse/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512102218_add_check_path",
		Migrate: func(tx *gorm.DB) error {
			check := &models.Check{}

			// Add path column if it doesn't exist
			if !tx.Migrator().HasColumn(check, "path") {
				if err := tx.Migrator().AddColumn(check, "path"); err != nil {
					return err
				}
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			check := &models.Check{}

			// Remove path column
			if tx.Migrator().HasColumn(check, "path") {
				if err := tx.Migrator().DropColumn(check, "path"); err != nil {
					return err
				}
			}

			return nil
		},
	})
}
