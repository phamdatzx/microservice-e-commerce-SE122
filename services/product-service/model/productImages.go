package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductImages struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProductID uuid.UUID      `gorm:"type:uuid;not null;index" json:"product_id"`
	Image     string         `gorm:"not null" json:"image"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationship
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
