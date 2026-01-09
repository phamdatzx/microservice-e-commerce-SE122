package stripeclient

import (
	"strings"
	"time"

	"order-service/model"

	"github.com/stripe/stripe-go/v84"
)

/*
	MapStripePaymentIntentToOrder
	- Chỉ map trạng thái + metadata
	- KHÔNG ghi DB
	- KHÔNG xử lý nghiệp vụ
*/
func MapStripePaymentIntentToOrder(
	order *model.Order,
	pi *stripe.PaymentIntent,
) *model.Order {

	if order == nil || pi == nil {
		return order
	}

	order.PaymentMethod = "STRIPE"
	order.PaymentStatus = mapStripePaymentStatus(pi.Status)
	order.UpdatedAt = time.Now()

	// Nếu thanh toán thành công → cập nhật order status tổng
	if order.PaymentStatus == "PAID" {
		order.Status = "completed"
	}

	// Nếu thất bại
	if order.PaymentStatus == "FAILED" {
		order.Status = "payment_failed"
	}

	return order
}

func mapStripePaymentStatus(
	status stripe.PaymentIntentStatus,
) string {

	switch status {

	case stripe.PaymentIntentStatusSucceeded:
		return "PAID"

	case stripe.PaymentIntentStatusCanceled:
		return "FAILED"

	case stripe.PaymentIntentStatusProcessing,
		stripe.PaymentIntentStatusRequiresAction,
		stripe.PaymentIntentStatusRequiresConfirmation,
		stripe.PaymentIntentStatusRequiresPaymentMethod:
		return "PENDING"

	default:
		return strings.ToUpper(string(status))
	}
}
