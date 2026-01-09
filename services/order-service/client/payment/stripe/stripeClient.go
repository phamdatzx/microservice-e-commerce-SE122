package stripeclient

import (
	"context"
	"fmt"

	"order-service/config"
	"order-service/model"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
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

func (s *StripeClient) CreatePayment(order *model.Order, successURL, cancelURL string) (string, error) {
	var lineItems []*stripe.CheckoutSessionLineItemParams
	
	for _, item := range order.Items {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("usd"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(item.ProductName + " - " + item.VariantName),
				},
				UnitAmount: stripe.Int64(int64(item.Price * 100)),
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Currency: stripe.String("usd"),
		LineItems:  lineItems,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Metadata: map[string]string{
			"order_id": order.ID,
		},
	}


	sess, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	return sess.URL, nil
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
