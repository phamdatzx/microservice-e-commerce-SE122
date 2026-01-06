package repository

import (
	"user-service/model"

	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(address *model.Address) error
	GetByID(id string) (*model.Address, error)
	GetByUserID(userID string) ([]model.Address, error)
	Update(address *model.Address) error
	Delete(id string) error
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(address *model.Address) error {
	return r.db.Create(address).Error
}

func (r *addressRepository) GetByID(id string) (*model.Address, error) {
	var address model.Address
	err := r.db.First(&address, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) GetByUserID(userID string) ([]model.Address, error) {
	var addresses []model.Address
	err := r.db.Where("user_id = ?", userID).Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *addressRepository) Update(address *model.Address) error {
	return r.db.Save(address).Error
}

func (r *addressRepository) Delete(id string) error {
	return r.db.Delete(&model.Address{}, "id = ?", id).Error
}
