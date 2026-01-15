package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserFollow struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	UserID   uuid.UUID `gorm:"type:uuid;column:user_id;not null"`
	SellerID uuid.UUID `gorm:"type:uuid;column:seller_id;not null"`
	User     User      `gorm:"foreignKey:UserID;references:ID"`
	Seller   User      `gorm:"foreignKey:SellerID;references:ID"`
}

// Hook tự động sinh UUID trước khi tạo record
func (u *UserFollow) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

