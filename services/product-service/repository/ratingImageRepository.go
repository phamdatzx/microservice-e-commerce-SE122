package repository

import (
	"product-service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RatingImageRepository interface {
	Create(ratingImage *model.RatingImage) error
	FindByID(id uuid.UUID) (*model.RatingImage, error)
	FindByRatingID(ratingID uuid.UUID) ([]model.RatingImage, error)
	FindAll() ([]model.RatingImage, error)
	Update(ratingImage *model.RatingImage) error
	Delete(id uuid.UUID) error
}

type ratingImageRepository struct {
	db *gorm.DB
}

func NewRatingImageRepository(db *gorm.DB) RatingImageRepository {
	return &ratingImageRepository{db: db}
}

func (r *ratingImageRepository) Create(ratingImage *model.RatingImage) error {
	return r.db.Create(ratingImage).Error
}

func (r *ratingImageRepository) FindByID(id uuid.UUID) (*model.RatingImage, error) {
	var ratingImage model.RatingImage
	err := r.db.First(&ratingImage, "id = ?", id).Error
	return &ratingImage, err
}

func (r *ratingImageRepository) FindByRatingID(ratingID uuid.UUID) ([]model.RatingImage, error) {
	var ratingImages []model.RatingImage
	err := r.db.Where("rating_id = ?", ratingID).Find(&ratingImages).Error
	return ratingImages, err
}

func (r *ratingImageRepository) FindAll() ([]model.RatingImage, error) {
	var ratingImages []model.RatingImage
	err := r.db.Find(&ratingImages).Error
	return ratingImages, err
}

func (r *ratingImageRepository) Update(ratingImage *model.RatingImage) error {
	return r.db.Save(ratingImage).Error
}

func (r *ratingImageRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.RatingImage{}, "id = ?", id).Error
}
