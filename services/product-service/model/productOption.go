package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductOption struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProductID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"product_id"`
	Name              string         `gorm:"not null" json:"name"`
	Price             float64        `gorm:"not null" json:"price"`
	QuantityAvailable int            `gorm:"not null;default:0" json:"quantity_available"`
	SoldCount         int            `gorm:"default:0" json:"sold_count"`
	Image             string         `json:"image"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationship
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
