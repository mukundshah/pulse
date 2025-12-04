package main

import (
	"flag"
	"log"
	"pulse/internal/db"

	"github.com/go-gormigrate/gormigrate/v2"

	"pulse/internal/config"
	"pulse/internal/models/migrations"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Usage: migrate [up|down] [id]")
	}

	direction := args[0]
	var id string
	if len(args) > 1 {
		id = args[1]
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	m := gormigrate.New(
		db,
		&gormigrate.Options{
			TableName:                 "__gorm_migrations__",
			IDColumnName:              "id",
			IDColumnSize:              255,
			UseTransaction:            false,
			ValidateUnknownMigrations: false,
		},
		migrations.GetAllMigrations(),
	)

	switch direction {
	case "up":
		if id != "" {
			if err := m.MigrateTo(id); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
			log.Printf("Migrated to: %s", id)
		} else {
			if err := m.Migrate(); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
			log.Println("Migration completed successfully")
		}
	case "down":
		if id != "" {
			if err := m.RollbackTo(id); err != nil {
				log.Fatalf("Rollback failed: %v", err)
			}
			log.Printf("Rolled back to: %s", id)
		} else {
			if err := m.RollbackLast(); err != nil {
				log.Fatalf("Rollback failed: %v", err)
			}
			log.Println("Rolled back last migration")
		}
	default:
		log.Fatalf("Invalid direction: %s. Use 'up' or 'down'", direction)
	}
}
