package dto

import "time"

type VoucherRequest struct {
	Code                   string    `json:"code" validate:"required"`
	Name                   string    `json:"name" validate:"required"`
	Description            string    `json:"description"`
	DiscountType           string    `json:"discount_type" validate:"required,oneof=FIXED PERCENTAGE"`
	DiscountValue          int       `json:"discount_value" validate:"required,min=0"`
	MaxDiscountValue       *int      `json:"max_discount_value"`
	MinOrderValue          int       `json:"min_order_value" validate:"min=0"`
	ApplyScope             string    `json:"apply_scope" validate:"required,oneof=ALL CATEGORY"`
	ApplySellerCategoryIds []string  `json:"apply_seller_category_ids"`
	TotalQuantity          int       `json:"total_quantity" validate:"required,min=1"`
	UsageLimitPerUser      int       `json:"usage_limit_per_user" validate:"required,min=1"`
	StartTime              time.Time `json:"start_time" validate:"required"`
	EndTime                time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	Status                 string    `json:"status" validate:"oneof=ACTIVE INACTIVE"`
}

type VoucherResponse struct {
	ID                     string    `json:"id"`
	SellerID               string    `json:"seller_id"`
	Code                   string    `json:"code"`
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	DiscountType           string    `json:"discount_type"`
	DiscountValue          int       `json:"discount_value"`
	MaxDiscountValue       *int      `json:"max_discount_value"`
	MinOrderValue          int       `json:"min_order_value"`
	ApplyScope             string    `json:"apply_scope"`
	ApplySellerCategoryIds []string  `json:"apply_seller_category_ids"`
	TotalQuantity          int       `json:"total_quantity"`
	UsedQuantity           int       `json:"used_quantity"`
	UsageLimitPerUser      int       `json:"usage_limit_per_user"`
	StartTime              time.Time `json:"start_time"`
	EndTime                time.Time `json:"end_time"`
	Status                 string    `json:"status"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}
