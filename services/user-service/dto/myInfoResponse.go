package dto

type MyInfoResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
	IsVerify bool   `json:"is_verify"`
	IsBanned bool   `json:"is_banned"`
}
