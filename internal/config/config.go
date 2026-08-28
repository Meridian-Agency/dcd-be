package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string   `env:"APP_ENV" envDefault:"production"`
	AllowedOrigins  []string `env:"ALLOWED_ORIGINS" envSeparator:"," envDefault:"https://demon-car-detailing-phi.vercel.app,http://localhost:3000"`
	Port            int      `env:"PORT" envDefault:"8080"`
	DBHost          string   `env:"DB_HOST,required"`
	DBPort          int      `env:"DB_PORT" envDefault:"5432"`
	DBDatabase      string   `env:"DB_DATABASE,required"`
	DBUsername      string   `env:"DB_USERNAME,required"`
	DBPassword      string   `env:"DB_PASSWORD,required" json:"-"`
	DBSchema        string   `env:"DB_SCHEMA" envDefault:"public"`
	DBSSLMode       string   `env:"DB_SSLMODE" envDefault:"disable"`
	AdminAPIKey     string   `env:"ADMIN_API_KEY" envDefault:"dcd-admin-secret"`
}

func Load() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load()
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
