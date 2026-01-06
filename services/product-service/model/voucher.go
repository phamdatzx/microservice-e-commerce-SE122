package model

import (
	"time"

	"github.com/google/uuid"
)

type Voucher struct {
	ID                     string    `bson:"_id" json:"id"`
	SellerID               string    `bson:"seller_id" json:"seller_id"`
	Code                   string    `bson:"code" json:"code"`
	Name                   string    `bson:"name" json:"name"`
	Description            string    `bson:"description" json:"description"`
	DiscountType           string    `bson:"discount_type" json:"discount_type"` // FIXED, PERCENTAGE
	DiscountValue          int       `bson:"discount_value" json:"discount_value"`
	MaxDiscountValue       *int      `bson:"max_discount_value" json:"max_discount_value"`
	MinOrderValue          int       `bson:"min_order_value" json:"min_order_value"`
	ApplyScope             string    `bson:"apply_scope" json:"apply_scope"` // ALL, CATEGORY
	ApplySellerCategoryIds []string  `bson:"apply_seller_category_ids" json:"apply_seller_category_ids"`
	TotalQuantity          int       `bson:"total_quantity" json:"total_quantity"`
	UsedQuantity           int       `bson:"used_quantity" json:"used_quantity"`
	UsageLimitPerUser      int       `bson:"usage_limit_per_user" json:"usage_limit_per_user"`
	StartTime              time.Time `bson:"start_time" json:"start_time"`
	EndTime                time.Time `bson:"end_time" json:"end_time"`
	Status                 string    `bson:"status" json:"status"` // ACTIVE, INACTIVE, EXPIRED
	CreatedAt              time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt              time.Time `bson:"updated_at" json:"updated_at"`
}

func (v *Voucher) BeforeCreate() {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	if v.CreatedAt.IsZero() {
		v.CreatedAt = time.Now()
	}
	if v.UpdatedAt.IsZero() {
		v.UpdatedAt = time.Now()
	}
	if v.UsedQuantity == 0 {
		v.UsedQuantity = 0
	}
	if v.Status == "" {
		v.Status = "ACTIVE"
	}
}

func (v *Voucher) BeforeUpdate() {
	v.UpdatedAt = time.Now()
}
