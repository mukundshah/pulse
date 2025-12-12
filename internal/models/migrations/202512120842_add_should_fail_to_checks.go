package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512120842_add_should_fail_to_checks",
		Migrate: func(tx *gorm.DB) error {
			// Add should_fail column
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN should_fail BOOLEAN NOT NULL DEFAULT false").Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Remove should_fail column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN should_fail").Error; err != nil {
				return err
			}

			return nil
		},
	})
}
