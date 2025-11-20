package repository

import (
	"product-service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductOptionRepository interface {
	Create(productOption *model.ProductOption) error
	FindByID(id uuid.UUID) (*model.ProductOption, error)
	FindByProductID(productID uuid.UUID) ([]model.ProductOption, error)
	FindAll() ([]model.ProductOption, error)
	Update(productOption *model.ProductOption) error
	Delete(id uuid.UUID) error
}

type productOptionRepository struct {
	db *gorm.DB
}

func NewProductOptionRepository(db *gorm.DB) ProductOptionRepository {
	return &productOptionRepository{db: db}
}

func (r *productOptionRepository) Create(productOption *model.ProductOption) error {
	return r.db.Create(productOption).Error
}

func (r *productOptionRepository) FindByID(id uuid.UUID) (*model.ProductOption, error) {
	var productOption model.ProductOption
	err := r.db.First(&productOption, "id = ?", id).Error
	return &productOption, err
}

func (r *productOptionRepository) FindByProductID(productID uuid.UUID) ([]model.ProductOption, error) {
	var productOptions []model.ProductOption
	err := r.db.Where("product_id = ?", productID).Find(&productOptions).Error
	return productOptions, err
}

func (r *productOptionRepository) FindAll() ([]model.ProductOption, error) {
	var productOptions []model.ProductOption
	err := r.db.Find(&productOptions).Error
	return productOptions, err
}

func (r *productOptionRepository) Update(productOption *model.ProductOption) error {
	return r.db.Save(productOption).Error
}

func (r *productOptionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.ProductOption{}, "id = ?", id).Error
}
