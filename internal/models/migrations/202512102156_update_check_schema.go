package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512102156_update_check_schema",
		Migrate: func(tx *gorm.DB) error {
			// Rename url column to host
			if err := tx.Exec("ALTER TABLE checks RENAME COLUMN url TO host").Error; err != nil {
				return err
			}

			// Add host column (in case url column didn't exist)
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN host VARCHAR NOT NULL DEFAULT ''").Error; err != nil {
				return err
			}

			// Add port column
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN port INTEGER DEFAULT 80").Error; err != nil {
				return err
			}

			// Add secure column
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN secure BOOLEAN DEFAULT false").Error; err != nil {
				return err
			}

			// Add query_params column
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN query_params JSONB").Error; err != nil {
				return err
			}

			// Add body column
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN body JSONB").Error; err != nil {
				return err
			}

			// Add ip_version column
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN ip_version VARCHAR NOT NULL DEFAULT 'ipv4'").Error; err != nil {
				return err
			}

			// Add ssl_verification column
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN ssl_verification BOOLEAN DEFAULT true").Error; err != nil {
				return err
			}

			// Add follow_redirects column
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN follow_redirects BOOLEAN DEFAULT true").Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Remove follow_redirects column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN follow_redirects").Error; err != nil {
				return err
			}

			// Remove ssl_verification column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN ssl_verification").Error; err != nil {
				return err
			}

			// Remove ip_version column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN ip_version").Error; err != nil {
				return err
			}

			// Remove body column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN body").Error; err != nil {
				return err
			}

			// Remove query_params column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN query_params").Error; err != nil {
				return err
			}

			// Remove secure column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN secure").Error; err != nil {
				return err
			}

			// Remove port column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN port").Error; err != nil {
				return err
			}

			// Rename host back to url
			if err := tx.Exec("ALTER TABLE checks RENAME COLUMN host TO url").Error; err != nil {
				return err
			}

			return nil
		},
	})
}
