package dto

import "time"

type VoucherResponse struct {
	ID                     string    `json:"id"`
	Code                   string    `json:"code"`
	DiscountType           string    `json:"discount_type"`
	DiscountValue          int       `json:"discount_value"`
	MinOrderValue          int       `json:"min_order_value"`
	MaxDiscountValue       int       `json:"max_discount_value"`
	StartTime              time.Time `json:"start_time"`
	EndTime                time.Time `json:"end_time"`
	UsageLimit             int       `json:"usage_limit"`
	UsedCount              int       `json:"used_count"`
	Status                 string    `json:"status"`
	SellerID               string    `json:"seller_id"`
	ApplySellerCategoryIds []string  `json:"apply_seller_category_ids"`
	ApplyScope             string    `json:"apply_scope"`
}

type UseVoucherRequest struct {
	UserID    string `json:"user_id"`
	VoucherID string `json:"voucher_id"`
}

type UseVoucherResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	UsageID   string    `json:"usage_id,omitempty"`
	UsedAt    time.Time `json:"used_at,omitempty"`
	VoucherID string    `json:"voucher_id,omitempty"`
}
