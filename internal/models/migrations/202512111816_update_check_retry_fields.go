package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512111816_update_check_retry_fields",
		Migrate: func(tx *gorm.DB) error {
			// Add new interval column first (before dropping interval_seconds)
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN interval VARCHAR NOT NULL DEFAULT '10m'").Error; err != nil {
				return err
			}

			// Convert existing interval_seconds to interval format
			if err := tx.Exec(`
				UPDATE checks
				SET interval = CASE
					WHEN interval_seconds < 60 THEN interval_seconds::text || 's'
					WHEN interval_seconds < 3600 THEN (interval_seconds / 60)::text || 'm'
					ELSE (interval_seconds / 3600)::text || 'h'
				END
				WHERE interval_seconds IS NOT NULL
			`).Error; err != nil {
				return err
			}

			// Remove old columns (after data migration)
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN expected_status").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN should_fail").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN timeout_ms").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN interval_seconds").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN alert_threshold").Error; err != nil {
				return err
			}

			// Add threshold columns
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN degraded_threshold INTEGER NOT NULL DEFAULT 3000").Error; err != nil {
				return err
			}

			if err := tx.Exec("ALTER TABLE checks ADD COLUMN degraded_threshold_unit VARCHAR(2) DEFAULT 'ms'").Error; err != nil {
				return err
			}

			if err := tx.Exec("ALTER TABLE checks ADD COLUMN failed_threshold INTEGER NOT NULL DEFAULT 5000").Error; err != nil {
				return err
			}

			if err := tx.Exec("ALTER TABLE checks ADD COLUMN failed_threshold_unit VARCHAR(2) DEFAULT 'ms'").Error; err != nil {
				return err
			}

			// Add retry columns
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN retries VARCHAR(20) DEFAULT 'none'").Error; err != nil {
				return err
			}

			if err := tx.Exec("ALTER TABLE checks ADD COLUMN retries_count INTEGER").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN retries_delay INTEGER").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN retries_delay_unit VARCHAR(2)").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN retries_factor DOUBLE PRECISION").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN retries_jitter VARCHAR").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN retries_jitter_factor DOUBLE PRECISION").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN retries_max_delay INTEGER").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN retries_max_delay_unit VARCHAR(2)").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN retries_timeout INTEGER").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN retries_timeout_unit VARCHAR(2)").Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Remove new retry columns
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN retries_timeout_unit").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN retries_timeout").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN retries_max_delay_unit").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN retries_max_delay").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN retries_jitter_factor").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN retries_jitter").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN retries_factor").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN retries_delay_unit").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN retries_delay").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN retries_count").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN retries").Error; err != nil {
				return err
			}

			// Remove threshold columns
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN failed_threshold_unit").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN failed_threshold").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN degraded_threshold_unit").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN degraded_threshold").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks DROP COLUMN interval").Error; err != nil {
				return err
			}

			// Restore old columns (with defaults)
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN alert_threshold INTEGER").Error; err != nil {
				return err
			}
			if err := tx.Exec("UPDATE checks SET alert_threshold = 3 WHERE alert_threshold IS NULL OR alert_threshold = 0").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN interval_seconds INTEGER").Error; err != nil {
				return err
			}
			if err := tx.Exec("UPDATE checks SET interval_seconds = 600 WHERE interval_seconds IS NULL OR interval_seconds = 0").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN timeout_ms INTEGER").Error; err != nil {
				return err
			}
			if err := tx.Exec("UPDATE checks SET timeout_ms = 10000 WHERE timeout_ms IS NULL OR timeout_ms = 0").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN should_fail BOOLEAN").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ADD COLUMN expected_status INTEGER").Error; err != nil {
				return err
			}
			if err := tx.Exec("UPDATE checks SET expected_status = 200 WHERE expected_status IS NULL OR expected_status = 0").Error; err != nil {
				return err
			}

			return nil
		},
	})
}
