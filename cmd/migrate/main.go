package main

import (
	"context"
	"log"

	"dcd-be/internal/config"
	"dcd-be/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Connecting to database: %s:%d/%s", cfg.DBHost, cfg.DBPort, cfg.DBDatabase)
	dbService := database.New(cfg)

	log.Println("Running auto-migrations...")
	db := dbService.GetDB(context.Background())
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Seeding services...")
	if err := database.SeedServices(db); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	log.Println("Database migration and seeding completed successfully!")
}
