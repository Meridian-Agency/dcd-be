package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port            int    `env:"PORT" envDefault:"8080"`
	DBHost          string `env:"DB_HOST,required"`
	DBPort          int    `env:"DB_PORT" envDefault:"5432"`
	DBDatabase      string `env:"DB_DATABASE,required"`
	DBUsername      string `env:"DB_USERNAME,required"`
	DBPassword      string `env:"DB_PASSWORD,required" json:"-"`
	DBSchema        string `env:"DB_SCHEMA" envDefault:"public"`
	SMTPHost        string `env:"SMTP_HOST" envDefault:"mailpit"`
	SMTPPort        int    `env:"SMTP_PORT" envDefault:"1025"`
	S3Endpoint      string `env:"S3_ENDPOINT" envDefault:"http://minio:9000"`
	WhatsAppMockURL string `env:"WHATSAPP_MOCK_URL" envDefault:"http://whatsapp-mock:3000/api/messages"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
