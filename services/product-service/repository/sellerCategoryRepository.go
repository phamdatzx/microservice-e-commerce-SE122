package repository

import (
	"product-service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SellerCategoryRepository interface {
	Create(sellerCategory *model.SellerCategory) error
	FindByID(id uuid.UUID) (*model.SellerCategory, error)
	FindAll() ([]model.SellerCategory, error)
	Update(sellerCategory *model.SellerCategory) error
	Delete(id uuid.UUID) error
}

type sellerCategoryRepository struct {
	db *gorm.DB
}

func NewSellerCategoryRepository(db *gorm.DB) SellerCategoryRepository {
	return &sellerCategoryRepository{db: db}
}

func (r *sellerCategoryRepository) Create(sellerCategory *model.SellerCategory) error {
	return r.db.Create(sellerCategory).Error
}

func (r *sellerCategoryRepository) FindByID(id uuid.UUID) (*model.SellerCategory, error) {
	var sellerCategory model.SellerCategory
	err := r.db.First(&sellerCategory, "id = ?", id).Error
	return &sellerCategory, err
}

func (r *sellerCategoryRepository) FindAll() ([]model.SellerCategory, error) {
	var sellerCategories []model.SellerCategory
	err := r.db.Find(&sellerCategories).Error
	return sellerCategories, err
}

func (r *sellerCategoryRepository) Update(sellerCategory *model.SellerCategory) error {
	return r.db.Save(sellerCategory).Error
}

func (r *sellerCategoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.SellerCategory{}, "id = ?", id).Error
}
