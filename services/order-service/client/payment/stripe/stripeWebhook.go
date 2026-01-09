package stripeclient

import (
	"encoding/json"
	"io"
	"net/http"
	"order-service/service"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

type WebhookHandler struct {
	webhookSecret string
	orderService  service.OrderService
}

func NewWebhookHandler(secret string, orderService service.OrderService) *WebhookHandler {
	return &WebhookHandler{
		webhookSecret: secret,
		orderService:  orderService,
	}
}

func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
    payload, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    sig := r.Header.Get("Stripe-Signature")

    event, err := webhook.ConstructEvent(
        payload,
        sig,
        h.webhookSecret,
    )
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    switch event.Type {

	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		json.Unmarshal(event.Data.Raw, &pi)

		orderID := pi.Metadata["order_id"]
		if orderID != "" {
			_ = h.orderService.UpdateOrderPaymentStatus(r.Context(), orderID, "PAID")
		}

	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent
		json.Unmarshal(event.Data.Raw, &pi)

		orderID := pi.Metadata["order_id"]
		if orderID != "" {
			_ = h.orderService.UpdateOrderPaymentStatus(r.Context(), orderID, "FAILED")
		}

	}

    w.WriteHeader(http.StatusOK)
}
