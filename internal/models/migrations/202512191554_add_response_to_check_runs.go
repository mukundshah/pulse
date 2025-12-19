package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512191554_add_response_to_check_runs",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN response JSONB").Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN response").Error; err != nil {
				return err
			}
			return nil
		},
	})
}
