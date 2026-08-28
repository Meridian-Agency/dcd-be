package database

import (
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

type Migration struct {
	ID   string `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name"`
}

type MigrationFunc func(db *gorm.DB) error

type migrationStep struct {
	id   string
	run  MigrationFunc
}

var migrationList = []migrationStep{
	{
		id: "202608280001_initial_schema",
		run: func(db *gorm.DB) error {
			slog.Info("Running initial schema migration (ServicePackage & Reservation)...")
			return db.AutoMigrate(&ServicePackage{}, &Reservation{})
		},
	},
}

func RunMigrations(db *gorm.DB) error {
	// AutoMigrate the migrations registry table itself
	if err := db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to migrate migrations table: %w", err)
	}

	for _, step := range migrationList {
		var count int64
		if err := db.Model(&Migration{}).Where("id = ?", step.id).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check migration %s: %w", step.id, err)
		}

		if count == 0 {
			slog.Info("Applying database migration", slog.String("migration_id", step.id))
			err := db.Transaction(func(tx *gorm.DB) error {
				if err := step.run(tx); err != nil {
					return err
				}
				m := Migration{
					ID:   step.id,
					Name: step.id,
				}
				return tx.Create(&m).Error
			})
			if err != nil {
				return fmt.Errorf("migration %s failed: %w", step.id, err)
			}
			slog.Info("Successfully applied database migration", slog.String("migration_id", step.id))
		} else {
			slog.Debug("Migration already applied", slog.String("migration_id", step.id))
		}
	}

	return nil
}
