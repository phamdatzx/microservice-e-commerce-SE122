package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	SellerID    uuid.UUID      `gorm:"type:uuid;not null" json:"seller_id"`
	Rating      float64        `gorm:"default:0" json:"rating"`
	RateCount   int            `gorm:"default:0" json:"rate_count"`
	SoldCount   int            `gorm:"default:0" json:"sold_count"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Options []ProductOption `gorm:"foreignKey:ProductID" json:"options,omitempty"`
	Images  []ProductImages `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	Reports []ProductReport `gorm:"foreignKey:ProductID" json:"reports,omitempty"`
}
