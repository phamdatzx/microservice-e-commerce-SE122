package dto

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Role        string `json:"role"`
	UserId      string `json:"user_id"`
}
