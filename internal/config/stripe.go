package config

import (
	"log"

	"github.com/stripe/stripe-go/v78"
)

// InitStripe sets up the global configuration for the Stripe SDK.
func InitStripe(cfg *Config) {
	if cfg.StripeSecretKey == "" {
		log.Println("WARNING: STRIPE_SECRET_KEY is empty. Stripe services will not function properly.")
		return
	}
	stripe.Key = cfg.StripeSecretKey
	log.Println("Stripe SDK initialized successfully.")
}
