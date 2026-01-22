package dto

type InstantCheckoutItem struct {
	SellerID  string `json:"seller_id" binding:"required"`
	ProductID string `json:"product_id" binding:"required"`
	VariantID string `json:"variant_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type InstantCheckoutRequest struct {
	Items             []InstantCheckoutItem `json:"items" binding:"required,min=1"`
	VoucherID         string                `json:"voucher_id"`
	ShippingAddress   ShippingAddressDto    `json:"shipping_address" binding:"required"`
	PaymentMethod     string                `json:"payment_method" binding:"required,oneof=COD STRIPE"`
	DeliveryServiceID int                   `json:"delivery_service_id" binding:"required"`
}
