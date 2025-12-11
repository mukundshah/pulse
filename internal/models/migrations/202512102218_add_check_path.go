package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512102218_add_check_path",
		Migrate: func(tx *gorm.DB) error {
			// Add path column
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN path VARCHAR NOT NULL DEFAULT '/'").Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Remove path column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN path").Error; err != nil {
				return err
			}

			return nil
		},
	})
}
