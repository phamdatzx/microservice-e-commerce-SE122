package dto

import "github.com/google/uuid"

type UserFollowResponse struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	FullName string    `json:"full_name"`
	Image    string    `json:"image"`
}
