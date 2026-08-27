package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set test environment variables
	os.Setenv("PORT", "9999")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_DATABASE", "test")
	os.Setenv("DB_USERNAME", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("ADMIN_API_KEY", "test-secret")

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_DATABASE")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("ADMIN_API_KEY")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Port != 9999 {
		t.Errorf("expected Port to be 9999, got %d", cfg.Port)
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("expected DBHost to be 'localhost', got %s", cfg.DBHost)
	}
	if cfg.AdminAPIKey != "test-secret" {
		t.Errorf("expected AdminAPIKey to be 'test-secret', got %s", cfg.AdminAPIKey)
	}
}
