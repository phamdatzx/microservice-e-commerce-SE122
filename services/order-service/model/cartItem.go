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
	Product   CartProduct `bson:"product" json:"product"`
	Variant   CartVariant `bson:"variant" json:"variant"`
	Quantity  int         `bson:"quantity" json:"quantity"`
	CreatedAt time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time   `bson:"updated_at" json:"updated_at"`
}

type CartProduct struct {
	ID                string   `bson:"id" json:"id"`
	Name              string   `bson:"name" json:"name"`
	SellerID          string   `bson:"seller_id" json:"seller_id"`
	SellerCategoryIDs []string `bson:"seller_category_ids" json:"seller_category_ids"`
}

type CartVariant struct {
	ID      string            `bson:"id" json:"id"`
	SKU     string            `bson:"sku" json:"sku"`
	Options map[string]string `bson:"options" json:"options"`
	Price   int               `bson:"price" json:"price"`
	Stock   int               `bson:"stock" json:"stock"`
	Image   string            `bson:"image" json:"image"`
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
	if c.Product.ID == "" {
		return &ValidationError{Field: "product_id", Message: "product_id is required"}
	}
	if c.Variant.ID == "" {
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
