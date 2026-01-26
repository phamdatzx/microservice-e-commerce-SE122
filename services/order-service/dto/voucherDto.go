package dto

import "time"

type VoucherResponse struct {
	ID                     string    `json:"id"`
	Code                   string    `json:"code"`
	DiscountType           string    `json:"discount_type"`
	DiscountValue          int       `json:"discount_value"`
	MinOrderValue          int       `json:"min_order_value"`
	MaxDiscountValue       int       `json:"max_discount_value"`
	StartDate              time.Time `json:"start_date"`
	EndDate                time.Time `json:"end_date"`
	UsageLimit             int       `json:"usage_limit"`
	UsedCount              int       `json:"used_count"`
	Status                 string    `json:"status"`
	SellerID          string    `json:"seller_id"`
	ApplySellerCategoryIds []string  `json:"apply_seller_category_ids"`
	ApplyScope             string               `json:"apply_scope"`

}
