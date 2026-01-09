package payment

import (
	"context"
	"order-service/model"

	"github.com/stripe/stripe-go/v84"
)

type CreatePaymentResult struct {
	PaymentID    string
	ClientSecret string
	Status       string
}

type PaymentClient interface {
	CreatePayment(order *model.Order, successURL, cancelURL string) (string, error)
	CancelPayment(ctx context.Context, paymentID string) error
	RefundPayment(ctx context.Context, paymentID string) error
	ConstructEvent(payload []byte, signature string) (*stripe.Event, error)
}
