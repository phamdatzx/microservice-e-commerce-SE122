package model

import (
	"fmt"

	"github.com/google/uuid"
)

type SellerCategory struct {
	ID       string `bson:"_id" json:"id"`
	SellerID string `bson:"seller_id" json:"seller_id"`
	Name     string `bson:"name" json:"name"`
}

// Validate checks if the seller category has all required fields
func (c *SellerCategory) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("category name is required")
	}
	if c.SellerID == "" {
		return fmt.Errorf("seller_id is required")
	}
	return nil
}

func (c *SellerCategory) BeforeCreate() {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
}
