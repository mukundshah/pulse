package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512162242_add_ip_fields_to_check_runs",
		Migrate: func(tx *gorm.DB) error {
			// Add ip_version column
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN ip_version VARCHAR(10)").Error; err != nil {
				return err
			}

			// Add ip_address column
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN ip_address VARCHAR(45)").Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Remove ip_address column
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN ip_address").Error; err != nil {
				return err
			}

			// Remove ip_version column
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN ip_version").Error; err != nil {
				return err
			}

			return nil
		},
	})
}
