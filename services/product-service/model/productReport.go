package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductReport struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProductID uuid.UUID      `gorm:"type:uuid;not null;index" json:"product_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Type      string         `gorm:"not null" json:"type"`
	Content   string         `json:"content"`
	Status    string         `json:"status"`
	IsDeleted bool           `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationship
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
