package migrations

import (
	"pulse/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	regions := []models.Region{
		{Name: "Asia Pacific", Code: "apac"},
	}

	RegisterMigration(&gormigrate.Migration{
		ID: "202512042142_seed_regions",
		Migrate: func(tx *gorm.DB) error {
			for _, region := range regions {
				var existing models.Region
				err := tx.Where("code = ?", region.Code).First(&existing).Error
				if err == gorm.ErrRecordNotFound {
					if err := tx.Create(&region).Error; err != nil {
						return err
					}
				} else if err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			for _, region := range regions {
				if err := tx.Where("code = ?", region.Code).Delete(&models.Region{}).Error; err != nil {
					return err
				}
			}
			return nil
		},
	})
}
