package config

import (
	"os"
	"testing"
)

func TestLoad_Stripe(t *testing.T) {
	// Set test environment variables
	os.Setenv("PORT", "9999")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_DATABASE", "test")
	os.Setenv("DB_USERNAME", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("STRIPE_SECRET_KEY", "sk_test_123")
	os.Setenv("STRIPE_WEBHOOK_SECRET", "whsec_123")
	os.Setenv("STRIPE_CURRENCY", "gbp")

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_DATABASE")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("STRIPE_SECRET_KEY")
		os.Unsetenv("STRIPE_WEBHOOK_SECRET")
		os.Unsetenv("STRIPE_CURRENCY")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.StripeSecretKey != "sk_test_123" {
		t.Errorf("expected StripeSecretKey to be 'sk_test_123', got %s", cfg.StripeSecretKey)
	}
	if cfg.StripeWebhookSecret != "whsec_123" {
		t.Errorf("expected StripeWebhookSecret to be 'whsec_123', got %s", cfg.StripeWebhookSecret)
	}
	if cfg.StripeCurrency != "gbp" {
		t.Errorf("expected StripeCurrency to be 'gbp', got %s", cfg.StripeCurrency)
	}
}
