package model

import (
	"time"

	"github.com/google/uuid"
)

type StockReservation struct {
	ID        string    `bson:"_id" json:"id"`
	OrderID   string    `bson:"order_id" json:"order_id"`
	VariantID string    `bson:"variant_id" json:"variant_id"`
	Quantity  int       `bson:"quantity" json:"quantity"`
	Status    string    `bson:"status" json:"status"` // "RESERVED" or "RELEASED"
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// BeforeCreate generates a new UUID for the ID field if not set and initializes defaults
func (sr *StockReservation) BeforeCreate() {
	if sr.ID == "" {
		sr.ID = uuid.New().String()
	}
	if sr.CreatedAt.IsZero() {
		sr.CreatedAt = time.Now()
	}
	if sr.UpdatedAt.IsZero() {
		sr.UpdatedAt = time.Now()
	}
	// Set default status if not provided
	if sr.Status == "" {
		sr.Status = "RESERVED"
	}
}

// BeforeUpdate updates the UpdatedAt timestamp
func (sr *StockReservation) BeforeUpdate() {
	sr.UpdatedAt = time.Now()
}
