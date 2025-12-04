package migrations

import (
	"pulse/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(&gormigrate.Migration{
		ID: "202512042111_initial",
		Migrate: func(tx *gorm.DB) error {
			return tx.Migrator().CreateTable(
				&models.User{},
				&models.Project{},
				&models.Tag{},
				&models.Region{},
				&models.Check{},
				&models.CheckRun{},
				&models.ProjectMember{},
				&models.ProjectInvitation{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&models.User{},
				&models.Project{},
				&models.Tag{},
				&models.Region{},
				&models.Check{},
				&models.CheckRun{},
				&models.ProjectMember{},
				&models.ProjectInvitation{},
			)
		},
	})
}
