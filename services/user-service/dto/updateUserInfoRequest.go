package dto

type UpdateUserInfoRequest struct {
	Name  string `json:"name" validate:"min=2,max=100"`
	Phone string `json:"phone" validate:"max=15"`
	Email string `json:"email" validate:"email"`
}
