package dto

import "time"

// VariantDto represents variant details from product-service
type VariantDto struct {
	ID      string            `json:"id"`
	SKU     string            `json:"sku"`
	Options map[string]string `json:"options"`
	Price   int               `json:"price"`
	Stock   int               `json:"stock"`
	Image   string            `json:"image"`
}

// ProductVariantDto mirrors product-service's CartVariantDto
type ProductVariantDto struct {
	ProductName string     `json:"product_name"`
	Variant     VariantDto `json:"variant"`
}

// CartItemDetailDto represents enriched cart item with product info
type CartItemDetailDto struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	SellerID    string    `json:"seller_id"`
	ProductID   string    `json:"product_id"`
	VariantID   string    `json:"variant_id"`
	Quantity    int       `json:"quantity"`
	ProductName string    `json:"product_name"`
	Variant     VariantDto `json:"variant"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetCartItemsResponse represents the response for getting cart items
type GetCartItemsResponse struct {
	CartItems []CartItemDetailDto `json:"cart_items"`
}
