package seeders

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func SeedNews(db *sql.DB) error {
	// 1. Insert News
	var newsIDs []uuid.UUID
	rows, err := db.Query(`
		INSERT INTO news (title, content, status) VALUES
		('AI Breakthrough', 'A major breakthrough in AI has been announced...', 'PUBLISHED'),
		('New Health Tips', 'Discover the latest tips for a healthy lifestyle.', 'PUBLISHED'),
		('Market Trends 2026', 'Insights into the business market for the upcoming year.', 'DRAFT')
		ON CONFLICT DO NOTHING
		RETURNING id;
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return err
		}
		newsIDs = append(newsIDs, id)
	}

	// If no news was inserted (already exists), get existing IDs
	if len(newsIDs) == 0 {
		rows, err = db.Query("SELECT id FROM news LIMIT 3")
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var id uuid.UUID
			if err := rows.Scan(&id); err != nil {
				return err
			}
			newsIDs = append(newsIDs, id)
		}
	}

	// 2. Get Topic IDs
	var topicIDs []uuid.UUID
	tRows, err := db.Query("SELECT id FROM topic LIMIT 5")
	if err != nil {
		return err
	}
	defer tRows.Close()

	for tRows.Next() {
		var id uuid.UUID
		if err := tRows.Scan(&id); err != nil {
			return err
		}
		topicIDs = append(topicIDs, id)
	}

	// 3. Link News to Topics (dummy mapping)
	if len(newsIDs) > 0 && len(topicIDs) > 0 {
		for i, nID := range newsIDs {
			// Just pick a topic based on index for variety
			tID := topicIDs[i%len(topicIDs)]
			_, err = db.Exec(`
				INSERT INTO news_topics (news_id, topic_id) 
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING;
			`, nID, tID)
			if err != nil {
				fmt.Printf("Warning: failed to link news %s to topic %s: %v\n", nID, tID, err)
			}
		}
	}

	return nil
}
