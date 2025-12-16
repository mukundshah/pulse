package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512161107_convert_timestamps_to_timestamptz",
		Migrate: func(tx *gorm.DB) error {
			// Convert users table timestamps
			if err := tx.Exec("ALTER TABLE users ALTER COLUMN created_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE users ALTER COLUMN updated_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE users ALTER COLUMN deleted_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE users ALTER COLUMN last_login TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}

			// Convert projects table timestamps
			if err := tx.Exec("ALTER TABLE projects ALTER COLUMN created_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE projects ALTER COLUMN updated_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE projects ALTER COLUMN deleted_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}

			// Convert regions table timestamps
			if err := tx.Exec("ALTER TABLE regions ALTER COLUMN created_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE regions ALTER COLUMN updated_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE regions ALTER COLUMN deleted_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}

			// Convert tags table timestamps
			if err := tx.Exec("ALTER TABLE tags ALTER COLUMN created_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE tags ALTER COLUMN updated_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE tags ALTER COLUMN deleted_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}

			// Convert checks table timestamps
			if err := tx.Exec("ALTER TABLE checks ALTER COLUMN created_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ALTER COLUMN updated_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ALTER COLUMN deleted_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ALTER COLUMN last_run_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ALTER COLUMN next_run_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}

			// Convert check_runs table timestamps
			if err := tx.Exec("ALTER TABLE check_runs ALTER COLUMN created_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs ALTER COLUMN updated_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs ALTER COLUMN deleted_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}

			// Convert project_members table timestamps
			if err := tx.Exec("ALTER TABLE project_members ALTER COLUMN created_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE project_members ALTER COLUMN updated_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE project_members ALTER COLUMN deleted_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}

			// Convert project_invitations table timestamps
			if err := tx.Exec("ALTER TABLE project_invitations ALTER COLUMN created_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE project_invitations ALTER COLUMN updated_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE project_invitations ALTER COLUMN deleted_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}

			// Convert sessions table timestamps
			if err := tx.Exec("ALTER TABLE sessions ALTER COLUMN created_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE sessions ALTER COLUMN updated_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE sessions ALTER COLUMN deleted_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE sessions ALTER COLUMN expires_at TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE sessions ALTER COLUMN last_activity TYPE TIMESTAMPTZ").Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Convert sessions table timestamps back
			if err := tx.Exec("ALTER TABLE sessions ALTER COLUMN last_activity TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE sessions ALTER COLUMN expires_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE sessions ALTER COLUMN deleted_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE sessions ALTER COLUMN updated_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE sessions ALTER COLUMN created_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}

			// Convert project_invitations table timestamps back
			if err := tx.Exec("ALTER TABLE project_invitations ALTER COLUMN deleted_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE project_invitations ALTER COLUMN updated_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE project_invitations ALTER COLUMN created_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}

			// Convert project_members table timestamps back
			if err := tx.Exec("ALTER TABLE project_members ALTER COLUMN deleted_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE project_members ALTER COLUMN updated_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE project_members ALTER COLUMN created_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}

			// Convert check_runs table timestamps back
			if err := tx.Exec("ALTER TABLE check_runs ALTER COLUMN deleted_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs ALTER COLUMN updated_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE check_runs ALTER COLUMN created_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}

			// Convert checks table timestamps back
			if err := tx.Exec("ALTER TABLE checks ALTER COLUMN next_run_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ALTER COLUMN last_run_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ALTER COLUMN deleted_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ALTER COLUMN updated_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE checks ALTER COLUMN created_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}

			// Convert tags table timestamps back
			if err := tx.Exec("ALTER TABLE tags ALTER COLUMN deleted_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE tags ALTER COLUMN updated_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE tags ALTER COLUMN created_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}

			// Convert regions table timestamps back
			if err := tx.Exec("ALTER TABLE regions ALTER COLUMN deleted_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE regions ALTER COLUMN updated_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE regions ALTER COLUMN created_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}

			// Convert projects table timestamps back
			if err := tx.Exec("ALTER TABLE projects ALTER COLUMN deleted_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE projects ALTER COLUMN updated_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE projects ALTER COLUMN created_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}

			// Convert users table timestamps back
			if err := tx.Exec("ALTER TABLE users ALTER COLUMN last_login TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE users ALTER COLUMN deleted_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE users ALTER COLUMN updated_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE users ALTER COLUMN created_at TYPE TIMESTAMP").Error; err != nil {
				return err
			}

			return nil
		},
	})
}
