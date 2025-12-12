package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	regions := []struct {
		Name string
		Code string
	}{
		{Name: "Asia Pacific", Code: "apac"},
	}

	RegisterMigration(&gormigrate.Migration{
		ID: "202512042142_seed_regions",
		Migrate: func(tx *gorm.DB) error {
			for _, region := range regions {
				// Insert with ON CONFLICT DO NOTHING to avoid duplicates
				if err := tx.Exec(`
					INSERT INTO regions (name, code, created_at, updated_at)
					VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
					ON CONFLICT (code) DO NOTHING
				`, region.Name, region.Code).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			for _, region := range regions {
				if err := tx.Exec("DELETE FROM regions WHERE code = ?", region.Code).Error; err != nil {
					return err
				}
			}
			return nil
		},
	})
}
