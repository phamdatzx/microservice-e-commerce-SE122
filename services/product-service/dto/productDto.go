package dto

import "product-service/model"

// GetProductsQueryParams represents query parameters for getting products
type GetProductsQueryParams struct {
	Page          int    `form:"page"`
	Limit         int    `form:"limit"`
	Category      string `form:"category"`
	SellerCategory string `form:"seller_category"`
	Status        string `form:"status"`
	SortBy        string `form:"sort_by"`
	SortDirection string `form:"sort_direction"`
	Search        string `form:"search"`
}

// SetDefaults sets default values for query parameters
func (p *GetProductsQueryParams) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.SortDirection == "" {
		p.SortDirection = "asc"
	}
	// Validate sort direction
	if p.SortDirection != "asc" && p.SortDirection != "desc" {
		p.SortDirection = "asc"
	}
}

// GetSkip calculates the number of documents to skip for pagination
func (p *GetProductsQueryParams) GetSkip() int {
	return (p.Page - 1) * p.Limit
}

// PaginationMetadata represents pagination information
type PaginationMetadata struct {
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
	TotalItems   int64 `json:"total_items"`
	ItemsPerPage int   `json:"items_per_page"`
}

// PaginatedProductsResponse represents paginated products response
type PaginatedProductsResponse struct {
	Products   []model.Product    `json:"products"`
	Pagination PaginationMetadata `json:"pagination"`
}

// NewPaginatedProductsResponse creates a new paginated response
func NewPaginatedProductsResponse(products []model.Product, total int64, page, limit int) *PaginatedProductsResponse {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &PaginatedProductsResponse{
		Products: products,
		Pagination: PaginationMetadata{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalItems:   total,
			ItemsPerPage: limit,
		},
	}
}

// Stock Reservation DTOs

// ReserveStockItem represents a single item to reserve
type ReserveStockItem struct {
	VariantID string `json:"variant_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

// ReserveStockRequest represents the request to reserve stock
type ReserveStockRequest struct {
	OrderID string             `json:"order_id" binding:"required"`
	Items   []ReserveStockItem `json:"items" binding:"required,min=1"`
}

// ReleaseStockRequest represents the request to release stock
type ReleaseStockRequest struct {
	OrderID string `json:"order_id" binding:"required"`
}

// SearchProductsQueryParams represents query parameters for searching products
type SearchProductsQueryParams struct {
	Page          int      `form:"page"`
	Limit         int      `form:"limit"`
	MinPrice      *int     `form:"min_price"`
	MaxPrice      *int     `form:"max_price"`
	MinRating     *float64 `form:"min_rating"`
	MaxRating     *float64 `form:"max_rating"`
	CategoryIDs   string   `form:"category_ids"` // comma-separated list
	SearchQuery   string   `form:"search_query"`
	SortBy        string   `form:"sort_by"`        // rating, sold_count, price
	SortDirection string   `form:"sort_direction"` // asc, desc
}

// SetDefaults sets default values for search query parameters
func (p *SearchProductsQueryParams) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.SortDirection == "" {
		p.SortDirection = "asc"
	}
	// Validate sort direction
	if p.SortDirection != "asc" && p.SortDirection != "desc" {
		p.SortDirection = "asc"
	}
}

// GetSkip calculates the number of documents to skip for pagination
func (p *SearchProductsQueryParams) GetSkip() int {
	return (p.Page - 1) * p.Limit
}

// CheckProductStatusRequest represents the request to check product status
type CheckProductStatusRequest struct {
	UserID     string   `json:"user_id" binding:"required"`
	ProductIDs []string `json:"product_ids" binding:"required,min=1"`
}

// ProductStatusInfo represents status information for a single product
type ProductStatusInfo struct {
	ProductID  string `json:"product_id"`
	IsReported bool   `json:"is_reported"`
	IsRated    bool   `json:"is_rated"`
}

// CheckProductStatusResponse represents the response for checking product status
type CheckProductStatusResponse struct {
	Products []ProductStatusInfo `json:"products"`
}

