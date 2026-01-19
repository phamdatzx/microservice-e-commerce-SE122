package model

import (
	"time"

	"github.com/google/uuid"
)

type ViewHistory struct {
	ID        string    `bson:"_id" json:"id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	ProductID string    `bson:"product_id" json:"product_id"`
	ViewedAt  time.Time `bson:"viewed_at" json:"viewed_at"`
}

func (vh *ViewHistory) BeforeCreate() {
	if vh.ID == "" {
		vh.ID = uuid.New().String()
	}
	if vh.ViewedAt.IsZero() {
		vh.ViewedAt = time.Now()
	}
}
