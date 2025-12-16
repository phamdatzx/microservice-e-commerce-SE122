package model

import (
	"github.com/google/uuid"
)

type SellerCategory struct {
	ID        string    `bson:"_id" json:"id"`
	SellerID  string    `bson:"seller_id" json:"seller_id"`
	Name      string    `bson:"name" json:"name"`
}

func (c *SellerCategory) BeforeCreate() {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
}
