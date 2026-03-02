package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-app/internal/logging"
	"go-app/seeders"
	"log/slog"
)

func runSeeder(db *sql.DB, target string) error {
	ctx := context.Background()
	logging.LogInfo(ctx, "Seeding target", slog.String("target", target))

	switch target {
	case "all":
		if err := seeders.SeedUsers(db); err != nil {
			return fmt.Errorf("seeding users failed: %w", err)
		}
		if err := seeders.SeedTopics(db); err != nil {
			return fmt.Errorf("seeding topics failed: %w", err)
		}
		if err := seeders.SeedNews(db); err != nil {
			return fmt.Errorf("seeding news failed: %w", err)
		}
	case "users":
		if err := seeders.SeedUsers(db); err != nil {
			return fmt.Errorf("seeding users failed: %w", err)
		}
	case "topics":
		if err := seeders.SeedTopics(db); err != nil {
			return fmt.Errorf("seeding topics failed: %w", err)
		}
	case "news":
		if err := seeders.SeedNews(db); err != nil {
			return fmt.Errorf("seeding news failed: %w", err)
		}
	default:
		return errors.New("unknown seed target: " + target)
	}

	return nil
}
