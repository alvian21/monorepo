package seeders

import (
	"database/sql"
)

func SeedTopics(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO topic (name) VALUES
		('Technology'),
		('Health'),
		('Business'),
		('Sports'),
		('Entertainment')
		ON CONFLICT (name) DO NOTHING;
	`)
	return err
}
