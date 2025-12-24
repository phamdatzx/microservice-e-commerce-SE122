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
