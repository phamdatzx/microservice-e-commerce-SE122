package dto

type CheckoutRequest struct {
	CartItemIDs     []string           `json:"cart_item_ids" binding:"required,min=1"`
	VoucherID       string             `json:"voucher_id"`
	ShippingAddress ShippingAddressDto `json:"shipping_address" binding:"required"`
	PaymentMethod   string             `json:"payment_method" binding:"required,oneof=COD STRIPE"`
}

type ShippingAddressDto struct {
	FullName    string  `json:"full_name" binding:"required"`
	Phone       string  `json:"phone" binding:"required"`
	AddressLine string  `json:"address_line" binding:"required"`
	Ward        string  `json:"ward"`
	District    string  `json:"district"`
	Province    string  `json:"province" binding:"required"`
	Country     string  `json:"country" binding:"required"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type CheckoutResponse struct {
	OrderID      string  `json:"order_id"`
	TotalAmount  float64 `json:"total_amount"`
	Status       string  `json:"status"`
	ClientSecret string  `json:"client_secret,omitempty"`
	PaymentUrl   string  `json:"payment_url,omitempty"`
}