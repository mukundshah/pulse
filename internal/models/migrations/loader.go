package migrations

import (
	"sort"

	"github.com/go-gormigrate/gormigrate/v2"
)

var migrationRegistry []*gormigrate.Migration

// RegisterMigration is called by each migration's init() function to register itself.
// This allows migrations to be automatically discovered without manually adding them
// to a slice in migrate.go.
func RegisterMigration(migration *gormigrate.Migration) {
	migrationRegistry = append(migrationRegistry, migration)
}

// GetAllMigrations returns all registered migrations, sorted by their ID
// to ensure correct execution order.
func GetAllMigrations() []*gormigrate.Migration {
	// Sort migrations by ID to ensure correct execution order
	migrations := make([]*gormigrate.Migration, len(migrationRegistry))
	copy(migrations, migrationRegistry)

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].ID < migrations[j].ID
	})

	return migrations
}
