package model

import (
	"time"

	"github.com/google/uuid"
)

type VoucherUsage struct {
	ID        string    `bson:"_id" json:"id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	VoucherID string    `bson:"voucher_id" json:"voucher_id"`
	UsedAt    time.Time `bson:"used_at" json:"used_at"`
}

func (vu *VoucherUsage) BeforeCreate() {
	if vu.ID == "" {
		vu.ID = uuid.New().String()
	}
	if vu.UsedAt.IsZero() {
		vu.UsedAt = time.Now()
	}
}
