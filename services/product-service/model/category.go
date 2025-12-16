package model

import (
	"fmt"

	"github.com/google/uuid"
)

type Category struct {
	ID    string `bson:"_id" json:"id"`
	Name  string `bson:"name" json:"name"`
	Image string `bson:"image" json:"image"`
}

// Validate checks if the category has all required fields
func (c *Category) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("category name is required")
	}
	return nil
}

func (c *Category) BeforeCreate() {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
}
