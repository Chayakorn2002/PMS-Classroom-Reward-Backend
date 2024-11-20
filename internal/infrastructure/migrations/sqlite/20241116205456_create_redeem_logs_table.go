package migrations
import (
	"database/sql"
)

func init() {
	Migrations = append(Migrations, createRedeemLogsTable)
}
var createRedeemLogsTable = &Migration{
	Title: "20241116205456_create_redeem_logs_table.go",
	Up: func(db *sql.DB) error {
		_, err := db.Exec(`
			CREATE TABLE redeem_log (
				id TEXT PRIMARY KEY,
				serial TEXT NOT NULL,
				course_id TEXT NOT NULL,
				google_classroom_student_id TEXT NOT NULL,
				assignment_id TEXT NOT NULL,
				created_at DATETIME NOT NULL,
				created_by TEXT NOT NULL
			);
		`)
		if err != nil {
			return err
		}
		return nil
	},
		Down: func(db *sql.DB) error {
		_, err := db.Exec(`
			DROP TABLE redeem_log;
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
