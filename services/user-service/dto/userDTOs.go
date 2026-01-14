package dto

type UserResponse struct {
	ID       string           `json:"id"`
	Username string           `json:"username"`
	Name     string           `json:"name"`
	Phone    string           `json:"phone"`
	Email    string           `json:"email"`
	Image    string           `json:"image"`
	Address  *AddressResponse `json:"address,omitempty"` // Default address only
}

type SaleInfoResponse struct {
	FollowCount   int     `json:"follow_count"`
	RatingCount   int     `json:"rating_count"`
	RatingAverage float64 `json:"rating_average"`
	IsFollowing   bool    `json:"is_following"`
}

type SellerResponse struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Image    string            `json:"image"`
	Address  *AddressResponse  `json:"address,omitempty"`
	SaleInfo *SaleInfoResponse `json:"sale_info,omitempty"`
}
