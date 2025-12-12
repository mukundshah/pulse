package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512051146_remove_consecutive_fails",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN consecutive_fails").Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN consecutive_fails INTEGER NOT NULL DEFAULT 0").Error; err != nil {
				return err
			}
			return nil
		},
	})
}
