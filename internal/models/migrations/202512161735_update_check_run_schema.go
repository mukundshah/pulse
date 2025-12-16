package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512161735_update_check_run_schema",
		Migrate: func(tx *gorm.DB) error {
			// Add failure_reason column
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN failure_reason VARCHAR(50)").Error; err != nil {
				return err
			}

			// Rename response_status to response_status_code and make it nullable
			if err := tx.Exec("ALTER TABLE check_runs RENAME COLUMN response_status TO response_status_code").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs ALTER COLUMN response_status_code DROP NOT NULL").Error; err != nil {
				return err
			}

			// Add new timestamp columns (nullable first)
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN run_started_at TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN run_ended_at TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN request_started_at TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN first_byte_at TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN response_ended_at TIMESTAMPTZ").Error; err != nil {
				return err
			}

			// Pre-fill run_started_at and request_started_at with created_at for existing rows
			if err := tx.Exec("UPDATE check_runs SET run_started_at = created_at WHERE run_started_at IS NULL").Error; err != nil {
				return err
			}
			if err := tx.Exec("UPDATE check_runs SET request_started_at = created_at WHERE request_started_at IS NULL").Error; err != nil {
				return err
			}

			// Pre-fill response_ended_at with request_started_at + total_time_ms for existing rows
			if err := tx.Exec("UPDATE check_runs SET response_ended_at = request_started_at + (total_time_ms * INTERVAL '1 millisecond') WHERE response_ended_at IS NULL AND total_time_ms IS NOT NULL").Error; err != nil {
				return err
			}

			// Add NOT NULL constraint to run_started_at and request_started_at
			if err := tx.Exec("ALTER TABLE check_runs ALTER COLUMN run_started_at SET NOT NULL").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs ALTER COLUMN request_started_at SET NOT NULL").Error; err != nil {
				return err
			}

			// Add connection_reused column
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN connection_reused BOOLEAN NOT NULL DEFAULT FALSE").Error; err != nil {
				return err
			}

			// Add response_size_bytes column
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN response_size_bytes BIGINT").Error; err != nil {
				return err
			}

			// Remove total_time_ms column
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN total_time_ms").Error; err != nil {
				return err
			}

			// Remove remarks column
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN remarks").Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Restore remarks column
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN remarks TEXT").Error; err != nil {
				return err
			}

			// Restore total_time_ms column
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN total_time_ms INTEGER").Error; err != nil {
				return err
			}

			// Remove response_size_bytes column
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN response_size_bytes").Error; err != nil {
				return err
			}

			// Remove connection_reused column
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN connection_reused").Error; err != nil {
				return err
			}

			// Remove new timestamp columns
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN response_ended_at").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN first_byte_at").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN request_started_at").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN run_ended_at").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN run_started_at").Error; err != nil {
				return err
			}

			// Rename response_status_code back to response_status and make it NOT NULL again
			if err := tx.Exec("ALTER TABLE check_runs ALTER COLUMN response_status_code SET NOT NULL").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs RENAME COLUMN response_status_code TO response_status").Error; err != nil {
				return err
			}

			// Remove failure_reason column
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN failure_reason").Error; err != nil {
				return err
			}

			return nil
		},
	})
}
