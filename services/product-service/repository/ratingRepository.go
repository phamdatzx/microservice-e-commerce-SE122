package repository

import (
	"product-service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RatingRepository interface {
	Create(rating *model.Rating) error
	FindByID(id uuid.UUID) (*model.Rating, error)
	FindByUserID(userID uuid.UUID) ([]model.Rating, error)
	FindByParentID(parentID uuid.UUID) ([]model.Rating, error)
	FindAll() ([]model.Rating, error)
	Update(rating *model.Rating) error
	Delete(id uuid.UUID) error
}

type ratingRepository struct {
	db *gorm.DB
}

func NewRatingRepository(db *gorm.DB) RatingRepository {
	return &ratingRepository{db: db}
}

func (r *ratingRepository) Create(rating *model.Rating) error {
	return r.db.Create(rating).Error
}

func (r *ratingRepository) FindByID(id uuid.UUID) (*model.Rating, error) {
	var rating model.Rating
	err := r.db.Preload("Images").First(&rating, "id = ?", id).Error
	return &rating, err
}

func (r *ratingRepository) FindByUserID(userID uuid.UUID) ([]model.Rating, error) {
	var ratings []model.Rating
	err := r.db.Where("user_id = ?", userID).Preload("Images").Find(&ratings).Error
	return ratings, err
}

func (r *ratingRepository) FindByParentID(parentID uuid.UUID) ([]model.Rating, error) {
	var ratings []model.Rating
	err := r.db.Where("parent_id = ?", parentID).Preload("Images").Find(&ratings).Error
	return ratings, err
}

func (r *ratingRepository) FindAll() ([]model.Rating, error) {
	var ratings []model.Rating
	err := r.db.Preload("Images").Find(&ratings).Error
	return ratings, err
}

func (r *ratingRepository) Update(rating *model.Rating) error {
	return r.db.Save(rating).Error
}

func (r *ratingRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Rating{}, "id = ?", id).Error
}
