package model

import (
	"time"

	"github.com/google/uuid"
)

type SavedVoucher struct {
	ID             string    `bson:"_id" json:"id"`
	UserID         string    `bson:"user_id" json:"user_id"`
	VoucherID      string    `bson:"voucher_id" json:"voucher_id"`
	SavedAt        time.Time `bson:"saved_at" json:"saved_at"`
	UsedCount      int       `bson:"used_count" json:"used_count"`
	MaxUsesAllowed int       `bson:"max_uses_allowed" json:"max_uses_allowed"`
	IsDeleted      bool      `bson:"is_deleted" json:"is_deleted"`
}

func (sv *SavedVoucher) BeforeCreate() {
	if sv.ID == "" {
		sv.ID = uuid.New().String()
	}
	if sv.SavedAt.IsZero() {
		sv.SavedAt = time.Now()
	}
	// Initialize default values
	sv.UsedCount = 0
	sv.IsDeleted = false
}
