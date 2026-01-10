package dto

type CreatePaymentRequest struct {
	CartItemIDs []string `json:"cart_item_ids" binding:"required"`
}

type CreatePaymentResponse struct {
	PaymentUrl string `json:"payment_url"`
}
