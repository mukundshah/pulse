package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512102145_add_region_flag",
		Migrate: func(tx *gorm.DB) error {
			// Add flag column
			if err := tx.Exec("ALTER TABLE regions ADD COLUMN flag VARCHAR(10)").Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Remove flag column
			if err := tx.Exec("ALTER TABLE regions DROP COLUMN flag").Error; err != nil {
				return err
			}

			return nil
		},
	})
}
