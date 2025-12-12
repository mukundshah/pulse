package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512122351_add_dns_fields_to_checks",
		Migrate: func(tx *gorm.DB) error {
			// Add dns_record_type column
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN dns_record_type VARCHAR").Error; err != nil {
				return err
			}

			// Add dns_resolver column
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN dns_resolver VARCHAR").Error; err != nil {
				return err
			}

			// Add dns_resolver_port column
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN dns_resolver_port INTEGER").Error; err != nil {
				return err
			}

			// Add dns_resolver_protocol column
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN dns_resolver_protocol VARCHAR").Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Remove dns_resolver_protocol column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN dns_resolver_protocol").Error; err != nil {
				return err
			}

			// Remove dns_resolver_port column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN dns_resolver_port").Error; err != nil {
				return err
			}

			// Remove dns_resolver column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN dns_resolver").Error; err != nil {
				return err
			}

			// Remove dns_record_type column
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN dns_record_type").Error; err != nil {
				return err
			}

			return nil
		},
	})
}
