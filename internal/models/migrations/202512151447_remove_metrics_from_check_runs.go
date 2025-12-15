package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512151447_remove_metrics_from_check_runs",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN metrics").Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN metrics JSONB").Error; err != nil {
				return err
			}
			return nil
		},
	})
}
