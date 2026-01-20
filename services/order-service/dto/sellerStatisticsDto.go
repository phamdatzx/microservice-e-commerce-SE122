package dto

import "time"

// GetSellerStatisticsRequest contains query parameters for seller statistics
type GetSellerStatisticsRequest struct {
	From time.Time `form:"from" binding:"required"` // Start of time range
	To   time.Time `form:"to" binding:"required"`   // End of time range
	Type string    `form:"type"`                    // Optional: "day" or "month" for breakdown
}

// PeriodStatistics contains statistics for a specific period
type PeriodStatistics struct {
	Period     string  `json:"period"`      // Period identifier (e.g., "2026-01-20" for day, "2026-01" for month)
	OrderCount int     `json:"order_count"` // Number of orders in this period
	Revenue    float64 `json:"revenue"`     // Revenue for this period
}

// GetSellerStatisticsResponse contains seller statistics
type GetSellerStatisticsResponse struct {
	OrderCount   int                `json:"order_count"`             // Total number of orders in the time range
	TotalRevenue float64            `json:"total_revenue"`           // Sum of order totals for all orders
	Breakdown    []PeriodStatistics `json:"breakdown,omitempty"`     // Period-by-period breakdown when type is specified
}
