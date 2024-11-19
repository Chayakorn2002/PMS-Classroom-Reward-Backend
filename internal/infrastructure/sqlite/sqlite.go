package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Chayakorn2002/pms-classroom-backend/config"
	_ "github.com/mattn/go-sqlite3"
)

func OpenSQLiteDB(ctx context.Context) (*sql.DB, error) {
	config := config.ProvideConfig()

	db, err := sql.Open(
		"sqlite3",
		fmt.Sprintf(
			`file:%s`,
			config.Sqlite.Path,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.InfoContext(ctx, "🔌 Connected to SQLite database")

	return db, nil
}
