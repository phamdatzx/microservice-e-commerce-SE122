package config

import (
	"os"

	"github.com/stripe/stripe-go/v84"
)

type StripeConfig struct {
    SecretKey     string
    WebhookSecret string
    Currency      string
}

func NewStripeConfig() *StripeConfig {
    return &StripeConfig{
        SecretKey:     os.Getenv("STRIPE_SECRET_KEY"),
        WebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
        Currency:      os.Getenv("STRIPE_CURRENCY"),
    }
}

func (s *StripeConfig) Init() {
    stripe.Key = s.SecretKey
}
