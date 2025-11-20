package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Rating struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ParentID    *uuid.UUID     `gorm:"type:uuid;index" json:"parent_id"`
	Score       float64        `json:"score"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Parent  *Rating       `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies []Rating      `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	Images  []RatingImage `gorm:"foreignKey:RatingID" json:"images,omitempty"`
}
