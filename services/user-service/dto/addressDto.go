package dto

import "github.com/google/uuid"

type AddressRequest struct {
	FullName    string  `json:"full_name" validate:"required"`
	Phone       string  `json:"phone" validate:"required"`
	AddressLine string  `json:"address_line" validate:"required"`
	Ward        string  `json:"ward" validate:"required"`
	District    string  `json:"district" validate:"required"`
	Province    string  `json:"province" validate:"required"`
	Country     string  `json:"country" validate:"required"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Default     bool    `json:"default"`
}

type AddressResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	FullName    string    `json:"full_name"`
	Phone       string    `json:"phone"`
	AddressLine string    `json:"address_line"`
	Ward        string    `json:"ward"`
	District    string    `json:"district"`
	Province    string    `json:"province"`
	Country     string    `json:"country"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Default     bool      `json:"default"`
}
