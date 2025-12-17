package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512171301_add_alerts",
		Migrate: func(tx *gorm.DB) error {
			// Create alerts table
			if err := tx.Exec(`
				CREATE TABLE alerts (
					id UUID PRIMARY KEY DEFAULT uuidv7(),
					status VARCHAR(20) NOT NULL,
					run_id UUID NOT NULL,
					region_id UUID NOT NULL,
					project_id UUID NOT NULL,
					check_id UUID NOT NULL,
					created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMPTZ,
					FOREIGN KEY (run_id) REFERENCES check_runs(id),
					FOREIGN KEY (region_id) REFERENCES regions(id),
					FOREIGN KEY (project_id) REFERENCES projects(id),
					FOREIGN KEY (check_id) REFERENCES checks(id)
				)
			`).Error; err != nil {
				return err
			}

			// Create indexes for alerts
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_alerts_run_id ON alerts(run_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_alerts_region_id ON alerts(region_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_alerts_project_id ON alerts(project_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_alerts_check_id ON alerts(check_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_alerts_deleted_at ON alerts(deleted_at)`).Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Drop alerts table
			if err := tx.Exec(`DROP TABLE alerts`).Error; err != nil {
				return err
			}

			return nil
		},
	})
}
