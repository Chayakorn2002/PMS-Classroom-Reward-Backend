package migrations

import (
	"database/sql"
)

func init() {
	Migrations = append(Migrations, createUsersTable)
}

var createUsersTable = &Migration{
	Title: "20241030213658_create_users_table.go",
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
			CREATE TABLE users (
				id TEXT PRIMARY KEY,
				course_id TEXT NOT NULL,
				google_classroom_student_id TEXT NOT NULL,
				firstname TEXT NOT NULL,
				lastname TEXT NOT NULL,
				email TEXT NOT NULL,
				password TEXT NOT NULL,
				created_at DATETIME NOT NULL,
				created_by TEXT NOT NULL,
				updated_at DATETIME,
				updated_by TEXT
			);
		`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(db *sql.DB) error {
		_, err := db.Exec(`
			DROP TABLE users;	
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
