package dto

import "github.com/google/uuid"

// AdminUserResponse is used for admin to get all user info (excluding password)
type AdminUserResponse struct {
	ID        uuid.UUID           `json:"id"`
	Username  string              `json:"username"`
	Name      string              `json:"name"`
	Phone     string              `json:"phone"`
	Email     string              `json:"email"`
	Image     string              `json:"image"`
	Role      string              `json:"role"`
	IsActive  bool                `json:"is_active"`
	IsVerify  bool                `json:"is_verify"`
	IsBanned  bool                `json:"is_banned"`
	Addresses []AddressResponse   `json:"addresses,omitempty"`
	SaleInfo  *SaleInfoResponse   `json:"sale_info,omitempty"`
}

// GetAllUsersResponse is the paginated response for admin get all users
type GetAllUsersResponse struct {
	Users []AdminUserResponse `json:"users"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
}
