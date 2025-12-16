package model

import (
	"github.com/google/uuid"
)

type ProductImages struct {
	ID        string    `bson:"_id" json:"id"`
	URL       string    `bson:"url" json:"url"`
	Order     int       `bson:"order" json:"order"`
}

func (p *ProductImages) BeforeCreate() {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
}
