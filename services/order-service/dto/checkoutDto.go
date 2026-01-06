package dto

type CheckoutRequest struct {
	CartItemIDs     []string `json:"cart_item_ids" binding:"required,min=1"`
	VoucherID       string   `json:"voucher_id"`
	ShippingAddress struct {
		Name    string `json:"name" binding:"required"`
		Phone   string `json:"phone" binding:"required"`
		Address string `json:"address" binding:"required"`
		City    string `json:"city" binding:"required"`
	} `json:"shipping_address" binding:"required"`
}

type CheckoutResponse struct {
	OrderID     string  `json:"order_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
}
