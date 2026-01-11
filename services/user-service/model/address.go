package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	UserID      uuid.UUID `gorm:"type:uuid;column:user_id"`
	FullName    string    `json:"full_name" gorm:"column:full_name"`
	Phone       string    `json:"phone" gorm:"column:phone"`
	AddressLine string    `json:"address_line" gorm:"column:address_line"`
	Ward        string    `json:"ward" gorm:"column:ward"`
	District    string    `json:"district" gorm:"column:district"`
	Province    string    `json:"province" gorm:"column:province"`
	WardCode    string    `json:"ward_code" gorm:"column:ward_code"`
	DistrictCode string    `json:"district_code" gorm:"column:district_code"`
	ProvinceCode string    `json:"province_code" gorm:"column:province_code"`
	Country     string    `json:"country" gorm:"column:country"`
	Latitude    float64   `json:"latitude" gorm:"column:latitude"`
	Longitude   float64   `json:"longitude" gorm:"column:longitude"`
	Default     bool      `json:"default" gorm:"column:is_default;default:false"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	return
}