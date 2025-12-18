package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512181728_add_run_number_to_check_runs",
		Migrate: func(tx *gorm.DB) error {
			// Add run_number column (nullable initially for backfill)
			if err := tx.Exec("ALTER TABLE check_runs ADD COLUMN run_number INTEGER").Error; err != nil {
				return err
			}

			// Backfill existing check runs with sequential numbers per check
			// This uses a window function to assign row numbers per check_id
			if err := tx.Exec(`
				UPDATE check_runs
				SET run_number = subquery.row_num
				FROM (
					SELECT
						id,
						ROW_NUMBER() OVER (PARTITION BY check_id ORDER BY created_at ASC, id ASC) as row_num
					FROM check_runs
					WHERE deleted_at IS NULL
				) AS subquery
				WHERE check_runs.id = subquery.id
			`).Error; err != nil {
				return err
			}

			// Make run_number NOT NULL
			if err := tx.Exec("ALTER TABLE check_runs ALTER COLUMN run_number SET NOT NULL").Error; err != nil {
				return err
			}

			// Create unique constraint on (check_id, run_number) to ensure uniqueness per check
			if err := tx.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_check_runs_check_id_run_number ON check_runs(check_id, run_number) WHERE deleted_at IS NULL").Error; err != nil {
				return err
			}

			// Create trigger function that auto-increments run_number per check
			// This always sets the run_number, ignoring any value provided by the application
			if err := tx.Exec(`
				CREATE OR REPLACE FUNCTION set_check_run_number()
				RETURNS TRIGGER AS $$
				BEGIN
					-- Always calculate the next number for this check
					-- This ensures the counter is always managed at the database level
					SELECT COALESCE(MAX(run_number), 0) + 1
					INTO NEW.run_number
					FROM check_runs
					WHERE check_id = NEW.check_id
						AND deleted_at IS NULL;
					RETURN NEW;
				END;
				$$ LANGUAGE plpgsql;
			`).Error; err != nil {
				return err
			}

			// Create trigger that fires before insert
			if err := tx.Exec(`
				DROP TRIGGER IF EXISTS trigger_set_check_run_number ON check_runs;
				CREATE TRIGGER trigger_set_check_run_number
					BEFORE INSERT ON check_runs
					FOR EACH ROW
					EXECUTE FUNCTION set_check_run_number();
			`).Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Drop trigger
			if err := tx.Exec("DROP TRIGGER IF EXISTS trigger_set_check_run_number ON check_runs").Error; err != nil {
				return err
			}

			// Drop function
			if err := tx.Exec("DROP FUNCTION IF EXISTS set_check_run_number()").Error; err != nil {
				return err
			}

			// Drop unique index
			if err := tx.Exec("DROP INDEX IF EXISTS idx_check_runs_check_id_run_number").Error; err != nil {
				return err
			}

			// Drop column
			if err := tx.Exec("ALTER TABLE check_runs DROP COLUMN IF EXISTS run_number").Error; err != nil {
				return err
			}

			return nil
		},
	})
}
