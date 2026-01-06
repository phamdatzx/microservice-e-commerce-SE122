package dto

import (
	"order-service/model"
	"time"
)

// AddCartItemRequest represents the request body for adding a cart item
type AddCartItemRequest struct {
	SellerID  string `json:"seller_id" binding:"required"`
	ProductID string `json:"product_id" binding:"required"`
	VariantID string `json:"variant_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// UpdateCartItemQuantityRequest represents the request body for updating cart item quantity
type UpdateCartItemQuantityRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// CartItemResponse represents the response format for a cart item
type CartItemResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	SellerID  string    `json:"seller_id"`
	ProductID string    `json:"product_id"`
	VariantID string    `json:"variant_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToCartItemResponse converts a CartItem model to CartItemResponse DTO
func ToCartItemResponse(item *model.CartItem) *CartItemResponse {
	return &CartItemResponse{
		ID:        item.ID,
		UserID:    item.UserID,
		SellerID:  item.SellerID,
		ProductID: item.Product.ID,
		VariantID: item.Variant.ID,
		Quantity:  item.Quantity,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
