package repository

import (
	"product-service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductImagesRepository interface {
	Create(productImages *model.ProductImages) error
	FindByID(id uuid.UUID) (*model.ProductImages, error)
	FindByProductID(productID uuid.UUID) ([]model.ProductImages, error)
	FindAll() ([]model.ProductImages, error)
	Update(productImages *model.ProductImages) error
	Delete(id uuid.UUID) error
}

type productImagesRepository struct {
	db *gorm.DB
}

func NewProductImagesRepository(db *gorm.DB) ProductImagesRepository {
	return &productImagesRepository{db: db}
}

func (r *productImagesRepository) Create(productImages *model.ProductImages) error {
	return r.db.Create(productImages).Error
}

func (r *productImagesRepository) FindByID(id uuid.UUID) (*model.ProductImages, error) {
	var productImages model.ProductImages
	err := r.db.First(&productImages, "id = ?", id).Error
	return &productImages, err
}

func (r *productImagesRepository) FindByProductID(productID uuid.UUID) ([]model.ProductImages, error) {
	var productImages []model.ProductImages
	err := r.db.Where("product_id = ?", productID).Find(&productImages).Error
	return productImages, err
}

func (r *productImagesRepository) FindAll() ([]model.ProductImages, error) {
	var productImages []model.ProductImages
	err := r.db.Find(&productImages).Error
	return productImages, err
}

func (r *productImagesRepository) Update(productImages *model.ProductImages) error {
	return r.db.Save(productImages).Error
}

func (r *productImagesRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.ProductImages{}, "id = ?", id).Error
}
