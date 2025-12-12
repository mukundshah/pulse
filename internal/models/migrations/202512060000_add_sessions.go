package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512060000_add_sessions",
		Migrate: func(tx *gorm.DB) error {
			// Create sessions table
			if err := tx.Exec(`
				CREATE TABLE sessions (
					id UUID PRIMARY KEY DEFAULT uuidv4(),
					user_id UUID NOT NULL,
					jti VARCHAR(255) NOT NULL UNIQUE,
					user_agent TEXT,
					ip_address VARCHAR(45),
					is_active BOOLEAN NOT NULL DEFAULT true,
					expires_at TIMESTAMP NOT NULL,
					last_activity TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP,
					FOREIGN KEY (user_id) REFERENCES users(id)
				)
			`).Error; err != nil {
				return err
			}

			// Create indexes
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_sessions_is_active ON sessions(is_active)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_sessions_deleted_at ON sessions(deleted_at)`).Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Drop sessions table
			if err := tx.Exec(`DROP TABLE sessions`).Error; err != nil {
				return err
			}

			return nil
		},
	})
}
