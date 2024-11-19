package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"

	migration "github.com/Chayakorn2002/pms-classroom-backend/internal/infrastructure/migrations/sqlite"
)

func main() {
	ctx := context.Background()

	InitMigrationCommand(ctx)
}

func InitMigrationCommand(ctx context.Context) {
	makeMigrationFile := flag.Bool("migrate:make", false, "Create new migration file")
	migrationFilename := flag.String("name", "", "Migration file name")
	migrateUp := flag.Bool("migrate:up", false, "Up database schema migration")
	migrateDown := flag.Bool("migrate:down", false, "Down database schema migration")
	migrationDownStep := flag.Int("step", 0, "Down database schema migration step")
	migrationFlush := flag.Bool("migrate:flush", false, "Reset database schema migration")
	flag.Parse()

	switch {
	case *makeMigrationFile:
		handleMakeMigrationFile(ctx, migrationFilename)
	case *migrateUp:
		handleMigrateUp(ctx)
	case *migrateDown:
		handleMigrateDown(ctx, *migrationDownStep)
	case *migrationFlush:
		handleMigrateFlush(ctx)
	default:
		return
	}
}

func handleMakeMigrationFile(ctx context.Context, migrationFilename *string) {
	slog.InfoContext(ctx, "Creating new migration file ...")
	if *migrationFilename == "" {
		log.Fatalln("🚨 Migration file name is required. Please run go run main.go -migrate:make -name=<migration_name>")
	}
	migration.MakeMigration(ctx, migrationFilename)
	os.Exit(0)
}

func handleMigrateUp(ctx context.Context) {
	log.Println("Migrating up ...")
	migration.MigrateUp(ctx)
	os.Exit(0)
}

func handleMigrateDown(ctx context.Context, migrationDownStep int) {
	log.Println("Migrating down ...")
	if migrationDownStep == 0 {
		log.Fatalln("🚨 Migration down step is required. Please run `go run main.go -migrate:down -step=<migration_step>`")
	}
	migration.MigrateDown(ctx, migrationDownStep)
	os.Exit(0)
}

func handleMigrateFlush(ctx context.Context) {
	log.Println("Resetting migration ...")
	migration.MigrateFlush(ctx)
	os.Exit(0)
}
