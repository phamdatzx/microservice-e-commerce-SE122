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
	ID              string             `json:"id"`
	Status          string             `json:"status"`
	User            UserDto            `json:"user"`
	PaymentMethod   string             `json:"payment_method"`
	PaymentStatus   string             `json:"payment_status"`
	Seller          UserDto            `json:"seller"`
	Items           []OrderItemDto     `json:"items"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	Voucher         *OrderVoucherDto   `json:"voucher"`
	Total           float64            `json:"total"`
	Phone           string             `json:"phone"`
	ShippingAddress OrderAddressDto    `json:"shipping_address"`
	DeliveryCode    string             `json:"delivery_code"`
	ItemCount       int                `json:"item_count"` // Computed field
}

// UserDto represents user information in orders
type UserDto struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// OrderItemDto represents an order item
type OrderItemDto struct {
	ProductID   string `json:"product_id"`
	VariantID   string `json:"variant_id"`
	ProductName string `json:"product_name"`
	VariantName string `json:"variant_name"`
	SKU         string `json:"sku"`
	Price       int    `json:"price"`
	Image       string `json:"image"`
	Quantity    int    `json:"quantity"`
}

// OrderAddressDto represents shipping address
type OrderAddressDto struct {
	FullName    string  `json:"full_name"`
	Phone       string  `json:"phone"`
	AddressLine string  `json:"address_line"`
	Ward        string  `json:"ward"`
	District    string  `json:"district"`
	Province    string  `json:"province"`
	Country     string  `json:"country"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

// OrderVoucherDto represents voucher information
type OrderVoucherDto struct {
	Code                   string   `json:"code"`
	DiscountType           string   `json:"discount_type"`
	DiscountValue          int      `json:"discount_value"`
}

// GetOrdersBySellerRequest contains query parameters for seller orders with enhanced filtering
type GetOrdersBySellerRequest struct {
	Status        string `form:"status"`         // Filter by order status
	PaymentMethod string `form:"payment_method"` // Filter by payment method (COD, STRIPE)
	PaymentStatus string `form:"payment_status"` // Filter by payment status (PENDING, PAID, FAILED)
	Search        string `form:"search"`         // Search by order ID or phone
	Page          int    `form:"page"`           // Page number (default: 1)
	Limit         int    `form:"limit"`          // Items per page (default: 10, max: 100)
	SortBy        string `form:"sort_by"`        // Field to sort by: total, created_at
	SortOrder     string `form:"sort_order"`     // Sort order: asc, desc (default: desc)
}
