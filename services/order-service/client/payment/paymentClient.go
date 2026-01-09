package payment

import "context"

type CreatePaymentResult struct {
	PaymentID    string
	ClientSecret string
	Status       string
}

type PaymentClient interface {
	CreatePayment(ctx context.Context, amount int64, currency string, orderID string) (*CreatePaymentResult, error)
	CancelPayment(ctx context.Context, paymentID string) error
	RefundPayment(ctx context.Context, paymentID string) error
}
