package migrations

import (
	"pulse/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512102156_update_check_schema",
		Migrate: func(tx *gorm.DB) error {
			check := &models.Check{}

			// Rename url column to host if it exists
			if tx.Migrator().HasColumn(check, "url") && !tx.Migrator().HasColumn(check, "host") {
				if err := tx.Exec("ALTER TABLE checks RENAME COLUMN url TO host").Error; err != nil {
					return err
				}
			}

			// Add host column if it doesn't exist (in case url column didn't exist)
			if !tx.Migrator().HasColumn(check, "host") {
				if err := tx.Migrator().AddColumn(check, "host"); err != nil {
					return err
				}
			}
			if !tx.Migrator().HasColumn(check, "port") {
				if err := tx.Migrator().AddColumn(check, "port"); err != nil {
					return err
				}
			}
			if !tx.Migrator().HasColumn(check, "secure") {
				if err := tx.Migrator().AddColumn(check, "secure"); err != nil {
					return err
				}
			}
			if !tx.Migrator().HasColumn(check, "query_params") {
				if err := tx.Migrator().AddColumn(check, "query_params"); err != nil {
					return err
				}
			}
			if !tx.Migrator().HasColumn(check, "body") {
				if err := tx.Migrator().AddColumn(check, "body"); err != nil {
					return err
				}
			}
			if !tx.Migrator().HasColumn(check, "ip_version") {
				if err := tx.Migrator().AddColumn(check, "ip_version"); err != nil {
					return err
				}
			}
			if !tx.Migrator().HasColumn(check, "ssl_verification") {
				if err := tx.Migrator().AddColumn(check, "ssl_verification"); err != nil {
					return err
				}
			}
			if !tx.Migrator().HasColumn(check, "follow_redirects") {
				if err := tx.Migrator().AddColumn(check, "follow_redirects"); err != nil {
					return err
				}
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			check := &models.Check{}

			// Remove the new fields (in reverse order)
			if tx.Migrator().HasColumn(check, "follow_redirects") {
				if err := tx.Migrator().DropColumn(check, "follow_redirects"); err != nil {
					return err
				}
			}
			if tx.Migrator().HasColumn(check, "ssl_verification") {
				if err := tx.Migrator().DropColumn(check, "ssl_verification"); err != nil {
					return err
				}
			}
			if tx.Migrator().HasColumn(check, "ip_version") {
				if err := tx.Migrator().DropColumn(check, "ip_version"); err != nil {
					return err
				}
			}
			if tx.Migrator().HasColumn(check, "body") {
				if err := tx.Migrator().DropColumn(check, "body"); err != nil {
					return err
				}
			}
			if tx.Migrator().HasColumn(check, "query_params") {
				if err := tx.Migrator().DropColumn(check, "query_params"); err != nil {
					return err
				}
			}
			if tx.Migrator().HasColumn(check, "secure") {
				if err := tx.Migrator().DropColumn(check, "secure"); err != nil {
					return err
				}
			}
			if tx.Migrator().HasColumn(check, "port") {
				if err := tx.Migrator().DropColumn(check, "port"); err != nil {
					return err
				}
			}

			// Rename host back to url
			if tx.Migrator().HasColumn(check, "host") && !tx.Migrator().HasColumn(check, "url") {
				if err := tx.Exec("ALTER TABLE checks RENAME COLUMN host TO url").Error; err != nil {
					return err
				}
			}

			return nil
		},
	})
}
