package stripeclient

import (
	"context"
	"fmt"

	"order-service/client/payment"
	"order-service/config"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/paymentintent"
	"github.com/stripe/stripe-go/v84/refund"
)

type StripeClient struct {
	config *config.StripeConfig
}

func NewStripeClient(cfg *config.StripeConfig) *StripeClient {
	cfg.Init()
	return &StripeClient{config: cfg}
}

func (s *StripeClient) CreatePayment(ctx context.Context, amount int64, currency string, orderID string) (*payment.CreatePaymentResult, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		Metadata: map[string]string{
			"order_id": orderID,
		},
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return &payment.CreatePaymentResult{
		PaymentID:    pi.ID,
		ClientSecret: pi.ClientSecret,
		Status:       string(pi.Status),
	}, nil
}

func (s *StripeClient) CancelPayment(ctx context.Context, paymentID string) error {
	params := &stripe.PaymentIntentCancelParams{}
	_, err := paymentintent.Cancel(paymentID, params)
	if err != nil {
		return fmt.Errorf("failed to cancel payment intent: %w", err)
	}
	return nil
}

func (s *StripeClient) RefundPayment(ctx context.Context, paymentID string) error {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentID),
	}
	_, err := refund.New(params)
	if err != nil {
		return fmt.Errorf("failed to refund payment: %w", err)
	}
	return nil
}
