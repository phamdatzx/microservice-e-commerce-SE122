package dto

import "time"

// GetOrdersRequest contains query parameters for filtering and pagination
type GetOrdersRequest struct {
	Status    string `form:"status"`     // Filter by order status
	Page      int    `form:"page"`       // Page number (default: 1)
	Limit     int    `form:"limit"`      // Items per page (default: 10, max: 100)
	SortBy    string `form:"sort_by"`    // Field to sort by: total, created_at
	SortOrder string `form:"sort_order"` // Sort order: asc, desc (default: desc)
}

// GetOrdersResponse contains paginated order results
type GetOrdersResponse struct {
	Orders     []OrderDto `json:"orders"`
	TotalCount int64      `json:"total_count"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalPages int        `json:"total_pages"`
}

// OrderDto represents an order in the response
type OrderDto struct {
	ID              string        `json:"id"`
	Status          string        `json:"status"`
	PaymentMethod   string        `json:"payment_method"`
	PaymentStatus   string        `json:"payment_status"`
	Total           float64       `json:"total"`
	ItemCount       int           `json:"item_count"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}
