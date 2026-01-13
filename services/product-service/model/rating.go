package model

import (
	"time"
)

type Rating struct {
	ID        string `bson:"_id" json:"id"`
	ProductID string `bson:"product_id" json:"product_id"`
	VariantID string `bson:"variant_id" json:"variant_id"`

	UserID    string `bson:"user_id" json:"user_id"`
	Star      int    `bson:"star" json:"star"`
	Content   string `bson:"content,omitempty" json:"content,omitempty"`
	Image     string `bson:"image,omitempty" json:"image,omitempty"`

	RatingResponse []RatingResponse `bson:"rating_response,omitempty" json:"rating_response,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type RatingResponse struct {
	ID      string `bson:"_id" json:"id"`
	Content string `bson:"content,omitempty" json:"content,omitempty"`
}