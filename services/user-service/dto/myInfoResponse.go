package dto

type MyInfoResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Image    string `json:"image"`
	IsActive bool   `json:"is_active"`
	IsVerify bool   `json:"is_verify"`
	IsBanned bool   `json:"is_banned"`
}
