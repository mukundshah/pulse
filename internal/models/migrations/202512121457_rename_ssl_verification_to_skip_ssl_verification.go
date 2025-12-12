package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512121457_rename_ssl_verification_to_skip_ssl_verification",
		Migrate: func(tx *gorm.DB) error {
			// Add new column skip_ssl_verification with default false
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN skip_ssl_verification BOOLEAN NOT NULL DEFAULT false").Error; err != nil {
				return err
			}

			// Invert existing data: SSLVerification=true becomes SkipSSLVerification=false
			// SSLVerification=false becomes SkipSSLVerification=true
			// NULL values (which default to true) become false
			if err := tx.Exec("UPDATE checks SET skip_ssl_verification = NOT COALESCE(ssl_verification, true)").Error; err != nil {
				return err
			}

			// Drop the old column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN ssl_verification").Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Add back ssl_verification column with default true (nullable to match original)
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN ssl_verification BOOLEAN DEFAULT true").Error; err != nil {
				return err
			}

			// Invert existing data back: SkipSSLVerification=false becomes SSLVerification=true
			// SkipSSLVerification=true becomes SSLVerification=false
			if err := tx.Exec("UPDATE checks SET ssl_verification = NOT skip_ssl_verification").Error; err != nil {
				return err
			}

			// Drop the skip_ssl_verification column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN skip_ssl_verification").Error; err != nil {
				return err
			}

			return nil
		},
	})
}
