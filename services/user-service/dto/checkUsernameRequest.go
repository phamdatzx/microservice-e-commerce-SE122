package dto

type CheckUsernameRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}
