package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512042111_initial",
		Migrate: func(tx *gorm.DB) error {
			// Create users table
			if err := tx.Exec(`
				CREATE TABLE users (
					id UUID PRIMARY KEY DEFAULT uuidv4(),
					name VARCHAR NOT NULL,
					email VARCHAR NOT NULL UNIQUE,
					password_hash VARCHAR NOT NULL,
					email_verified BOOLEAN NOT NULL DEFAULT false,
					is_active BOOLEAN NOT NULL DEFAULT true,
					last_login TIMESTAMP,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				)
			`).Error; err != nil {
				return err
			}

			// Create indexes for users
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at)`).Error; err != nil {
				return err
			}

			// Create projects table
			if err := tx.Exec(`
				CREATE TABLE projects (
					id UUID PRIMARY KEY DEFAULT uuidv7(),
					name VARCHAR NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				)
			`).Error; err != nil {
				return err
			}

			// Create indexes for projects
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_projects_deleted_at ON projects(deleted_at)`).Error; err != nil {
				return err
			}

			// Create regions table
			if err := tx.Exec(`
				CREATE TABLE regions (
					id UUID PRIMARY KEY DEFAULT uuidv7(),
					name VARCHAR NOT NULL,
					code VARCHAR NOT NULL UNIQUE,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				)
			`).Error; err != nil {
				return err
			}

			// Create indexes for regions
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_regions_deleted_at ON regions(deleted_at)`).Error; err != nil {
				return err
			}

			// Create tags table
			if err := tx.Exec(`
				CREATE TABLE tags (
					id UUID PRIMARY KEY DEFAULT uuidv7(),
					name VARCHAR NOT NULL,
					project_id UUID NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (project_id) REFERENCES projects(id)
				)
			`).Error; err != nil {
				return err
			}

			// Create indexes for tags
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_tags_project_id ON tags(project_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_tags_deleted_at ON tags(deleted_at)`).Error; err != nil {
				return err
			}

			// Create checks table
			if err := tx.Exec(`
				CREATE TABLE checks (
					id UUID PRIMARY KEY DEFAULT uuidv7(),
					name VARCHAR NOT NULL,
					is_enabled BOOLEAN NOT NULL DEFAULT true,
					is_muted BOOLEAN NOT NULL DEFAULT false,
					type VARCHAR NOT NULL DEFAULT 'http',
					url VARCHAR NOT NULL,
					method VARCHAR NOT NULL DEFAULT 'GET',
					headers JSONB,
					playwright_script TEXT,
					assertions JSONB,
					expected_status INTEGER NOT NULL DEFAULT 200,
					should_fail BOOLEAN NOT NULL DEFAULT false,
					pre_script TEXT,
					post_script TEXT,
					timeout_ms INTEGER NOT NULL DEFAULT 10000,
					interval_seconds INTEGER NOT NULL,
					alert_threshold INTEGER NOT NULL DEFAULT 3,
					consecutive_fails INTEGER NOT NULL DEFAULT 0,
					last_status VARCHAR(20) NOT NULL DEFAULT 'unknown',
					last_run_at TIMESTAMP,
					next_run_at TIMESTAMP,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					project_id UUID NOT NULL,
					FOREIGN KEY (project_id) REFERENCES projects(id)
				)
			`).Error; err != nil {
				return err
			}

			// Create indexes for checks
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_checks_project_id ON checks(project_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_checks_deleted_at ON checks(deleted_at)`).Error; err != nil {
				return err
			}

			// Create check_runs table
			if err := tx.Exec(`
				CREATE TABLE check_runs (
					id UUID PRIMARY KEY DEFAULT uuidv7(),
					status VARCHAR(20) NOT NULL DEFAULT 'unknown',
					response_status INTEGER,
					assertion_results JSONB,
					playwright_report JSONB,
					network_timings JSONB,
					metrics JSONB,
					region_id UUID NOT NULL,
					check_id UUID NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (region_id) REFERENCES regions(id),
					FOREIGN KEY (check_id) REFERENCES checks(id)
				)
			`).Error; err != nil {
				return err
			}

			// Create indexes for check_runs
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_check_runs_region_id ON check_runs(region_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_check_runs_check_id ON check_runs(check_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_check_runs_deleted_at ON check_runs(deleted_at)`).Error; err != nil {
				return err
			}

			// Create project_members table
			if err := tx.Exec(`
				CREATE TABLE project_members (
					id UUID PRIMARY KEY DEFAULT uuidv4(),
					project_id UUID NOT NULL,
					user_id UUID NOT NULL,
					role VARCHAR(20) NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (project_id) REFERENCES projects(id),
					FOREIGN KEY (user_id) REFERENCES users(id)
				)
			`).Error; err != nil {
				return err
			}

			// Create indexes for project_members
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_project_members_project_id ON project_members(project_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_project_members_user_id ON project_members(user_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_project_members_deleted_at ON project_members(deleted_at)`).Error; err != nil {
				return err
			}

			// Create project_invitations table
			if err := tx.Exec(`
				CREATE TABLE project_invitations (
					id UUID PRIMARY KEY DEFAULT uuidv4(),
					project_id UUID NOT NULL,
					email VARCHAR NOT NULL,
					token VARCHAR NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (project_id) REFERENCES projects(id)
				)
			`).Error; err != nil {
				return err
			}

			// Create indexes for project_invitations
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_project_invitations_project_id ON project_invitations(project_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_project_invitations_deleted_at ON project_invitations(deleted_at)`).Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Drop tables in reverse order (respecting foreign key constraints)
			if err := tx.Exec(`DROP TABLE project_invitations`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`DROP TABLE project_members`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`DROP TABLE check_runs`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`DROP TABLE checks`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`DROP TABLE tags`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`DROP TABLE regions`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`DROP TABLE projects`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`DROP TABLE users`).Error; err != nil {
				return err
			}

			return nil
		},
	})
}
