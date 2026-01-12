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
