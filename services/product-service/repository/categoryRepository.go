package repository

import (
	"product-service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *model.Category) error
	FindByID(id uuid.UUID) (*model.Category, error)
	FindAll() ([]model.Category, error)
	Update(category *model.Category) error
	Delete(id uuid.UUID) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindByID(id uuid.UUID) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, "id = ?", id).Error
	return &category, err
}

func (r *categoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Category{}, "id = ?", id).Error
}
