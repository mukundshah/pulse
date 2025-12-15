package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512151520_add_total_time_ms_to_check_runs",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN total_time_ms INTEGER").Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN total_time_ms").Error; err != nil {
				return err
			}
			return nil
		},
	})
}
