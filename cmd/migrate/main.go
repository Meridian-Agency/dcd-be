package main

import (
	"context"
	"log/slog"
	"os"

	"dcd-be/internal/config"
	"dcd-be/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Connecting to database",
		slog.String("host", cfg.DBHost),
		slog.Int("port", cfg.DBPort),
		slog.String("database", cfg.DBDatabase),
	)
	dbService := database.New(cfg)

	slog.Info("Running database migrations...")
	db := dbService.GetDB(context.Background())
	if err := database.RunMigrations(db); err != nil {
		slog.Error("Migration failed", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Seeding services...")
	if err := database.SeedServices(db); err != nil {
		slog.Error("Seeding failed", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Database migration and seeding completed successfully!")
}
