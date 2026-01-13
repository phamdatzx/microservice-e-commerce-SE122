package model

import (
	"time"

	"github.com/google/uuid"
)

type Rating struct {
	ID        string `bson:"_id" json:"id"`
	ProductID string `bson:"product_id" json:"product_id"`
	VariantID string `bson:"variant_id" json:"variant_id"`

	User    User `bson:"user" json:"user"`
	Star      int    `bson:"star" json:"star"`
	Content   string `bson:"content,omitempty" json:"content,omitempty"`
	Image     string `bson:"image,omitempty" json:"image,omitempty"`

	RatingResponse []RatingResponse `bson:"rating_response,omitempty" json:"rating_response,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type User struct {
	ID string `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Image string `bson:"image" json:"image"`
	Phone string `bson:"phone" json:"phone"`
}

type RatingResponse struct {
	ID      string `bson:"_id" json:"id"`
	Content string `bson:"content,omitempty" json:"content,omitempty"`
}

func (r *Rating) BeforeCreate() {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
}

