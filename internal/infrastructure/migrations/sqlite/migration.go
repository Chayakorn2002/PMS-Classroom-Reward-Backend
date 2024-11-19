package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Chayakorn2002/pms-classroom-backend/internal/infrastructure/sqlite"
	array "github.com/Chayakorn2002/pms-classroom-backend/utils/array"
	"github.com/iancoleman/strcase"
)

type Migration struct {
	Title string
	Up    func(*sql.DB) error
	Down  func(*sql.DB) error
}

var Migrations []*Migration

func MakeMigration(ctx context.Context, migrationFilename *string) {
	formattedFilename := time.Now().Format("20060102150405") + "_" + *migrationFilename + ".go"
	filepath := "internal/infrastructure/migrations/sqlite/" + formattedFilename

	migrationVarName := strcase.ToLowerCamel(strings.ReplaceAll(*migrationFilename, "_", " "))

	formattedTemplate := strings.Replace(migrationTemplate, "<migration_name>", migrationVarName, -1)
	formattedTemplate = strings.Replace(formattedTemplate, "<filename>", formattedFilename, -1)

	_, err := os.Create(filepath)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	err = os.WriteFile(filepath, []byte(formattedTemplate), 0)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.InfoContext(ctx, fmt.Sprintf("📁 Migration file %s created successfully", formattedFilename))
}

const migrationTemplate = `package migrations
import (
	"database/sql"
)

func init() {
	Migrations = append(Migrations, <migration_name>)
}
var <migration_name> = &Migration{
	Title: "<filename>",
	Up: func(db *sql.DB) error {
		_, err := db.Exec(` + "`" + `
			// Write your migration query here
		` + "`" + `)
		if err != nil {
			return err
		}
		return nil
	},
		Down: func(db *sql.DB) error {
		_, err := db.Exec(` + "`" + `
			// Write your rollback query here
		` + "`" + `)
		if err != nil {
			return err
		}
		return nil
	},
}
`

func MigrateUp(ctx context.Context) {
	db, err := sqlite.OpenSQLiteDB(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}
	if db == nil {
		slog.ErrorContext(ctx, "🚨 SQLite database not found")
		return
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(255) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (title)
		);`,
	)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}

	migrations, err := db.Query(
		`SELECT title FROM migrations ORDER BY title;`,
	)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}
	defer migrations.Close()

	executedMigrations := make(map[string]interface{})
	for migrations.Next() {
		var title string
		if err := migrations.Scan(&title); err != nil {
			slog.ErrorContext(ctx, err.Error())
			return
		}
		executedMigrations[title] = nil
	}

	for _, migration := range Migrations {
		_, exists := executedMigrations[migration.Title]
		if !exists {
			slog.InfoContext(ctx, fmt.Sprintf("🚀 Migrating up %s ...", migration.Title))
			err := migration.Up(db)
			if err != nil {
				slog.Error(err.Error())
				return
			}

			_, err = db.Exec(
				`INSERT INTO migrations (title) VALUES ($1);`,
				migration.Title,
			)
			if err != nil {
				slog.Error(err.Error())
				return
			}
		}
	}

	slog.InfoContext(ctx, "🚀 Database is already up to date")
}

func MigrateDown(ctx context.Context, step int) {
	db, err := sqlite.OpenSQLiteDB(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}

	currentMigrations, err := db.Query(
		`SELECT title FROM migrations ORDER BY title DESC;`,
	)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}
	defer currentMigrations.Close()

	executedMigrations := []string{}
	for currentMigrations.Next() {
		var title string
		if err := currentMigrations.Scan(&title); err != nil {
			slog.Error(err.Error())
			return
		}
		executedMigrations = append(executedMigrations, title)
	}

	if step > len(executedMigrations) {
		step = len(executedMigrations)
	}

	for i := 0; i < step; i++ {
		for _, migration := range Migrations {
			if migration.Title == executedMigrations[i] {
				slog.InfoContext(ctx, fmt.Sprintf("🔙 Migrating down %s ...", migration.Title))
				err := migration.Down(db)
				if err != nil {
					slog.Error(err.Error())
					return
				}

				_, err = db.Exec(
					`DELETE FROM migrations WHERE title = $1;`,
					migration.Title,
				)
				if err != nil {
					slog.Error(err.Error())
					return
				}
				break
			}
		}
	}
}

var skippedTable = []string{
	"sqlite_sequence",
}

func MigrateFlush(ctx context.Context) {
	db, err := sqlite.OpenSQLiteDB(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}

	rows, err := db.Query(`
		SELECT name FROM sqlite_master WHERE type='table';
	`)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}
	defer rows.Close()

	var (
		tableName  string
		tableNames []string
	)

	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			slog.ErrorContext(ctx, err.Error())
			return
		}

		if !array.ContainAny([]string{tableName}, skippedTable) {
			tableNames = append(tableNames, tableName)
		}
	}

	for _, tableName := range tableNames {
		_, err = db.Exec(
			fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName),
		)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			return
		}
		slog.InfoContext(ctx, fmt.Sprintf("🔥 Flushing %s ...", tableName))
	}

	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}

	slog.InfoContext(ctx, "🔥 Migration flushed successfully")
}
