package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512062100_add_join_tables",
		Migrate: func(tx *gorm.DB) error {
			// Create check_tags join table
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS check_tags (
					check_id UUID NOT NULL,
					tag_id UUID NOT NULL,
					PRIMARY KEY (check_id, tag_id),
					FOREIGN KEY (check_id) REFERENCES checks(id) ON DELETE CASCADE,
					FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
				)
			`).Error; err != nil {
				return err
			}

			// Create check_regions join table
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS check_regions (
					check_id UUID NOT NULL,
					region_id UUID NOT NULL,
					PRIMARY KEY (check_id, region_id),
					FOREIGN KEY (check_id) REFERENCES checks(id) ON DELETE CASCADE,
					FOREIGN KEY (region_id) REFERENCES regions(id) ON DELETE CASCADE
				)
			`).Error; err != nil {
				return err
			}

			// Create indexes for better query performance
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_check_tags_check_id ON check_tags(check_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_check_tags_tag_id ON check_tags(tag_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_check_regions_check_id ON check_regions(check_id)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_check_regions_region_id ON check_regions(region_id)`).Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec(`DROP TABLE IF EXISTS check_regions`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`DROP TABLE IF EXISTS check_tags`).Error; err != nil {
				return err
			}
			return nil
		},
	})
}
