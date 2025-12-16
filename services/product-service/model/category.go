package model

import (
	"github.com/google/uuid"
)

type Category struct {
	ID        string    `bson:"_id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Image     string    `bson:"image" json:"image"`
}

func (c *Category) BeforeCreate() {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
}
