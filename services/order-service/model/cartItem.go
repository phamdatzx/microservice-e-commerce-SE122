package model

import (
	"time"

	"github.com/google/uuid"
)

// CartItem represents an item in a user's shopping cart
type CartItem struct {
	ID        string    `bson:"_id" json:"id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	SellerID  string    `bson:"seller_id" json:"seller_id"`
	ProductID string    `bson:"product_id" json:"product_id"`
	VariantID string    `bson:"variant_id" json:"variant_id"`
	Quantity  int       `bson:"quantity" json:"quantity"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// BeforeCreate generates a new UUID for the ID field if not set and initializes timestamps
func (c *CartItem) BeforeCreate() {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	if c.UpdatedAt.IsZero() {
		c.UpdatedAt = time.Now()
	}
}

// BeforeUpdate updates the UpdatedAt timestamp
func (c *CartItem) BeforeUpdate() {
	c.UpdatedAt = time.Now()
}

// Validate checks if the cart item has all required fields
func (c *CartItem) Validate() error {
	if c.UserID == "" {
		return &ValidationError{Field: "user_id", Message: "user_id is required"}
	}
	if c.SellerID == "" {
		return &ValidationError{Field: "seller_id", Message: "seller_id is required"}
	}
	if c.ProductID == "" {
		return &ValidationError{Field: "product_id", Message: "product_id is required"}
	}
	if c.VariantID == "" {
		return &ValidationError{Field: "variant_id", Message: "variant_id is required"}
	}
	if c.Quantity <= 0 {
		return &ValidationError{Field: "quantity", Message: "quantity must be greater than 0"}
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
