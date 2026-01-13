package service

import (
	"product-service/model"
	"product-service/repository"
	"time"

	"github.com/google/uuid"
)

type RatingService interface {
	CreateRating(rating *model.Rating) error
	GetRatingByID(id string) (*model.Rating, error)
	GetAllRatings() ([]model.Rating, error)
	GetRatingsByProductID(productID string, page, limit int) ([]model.Rating, int64, error)
	GetRatingsByUserID(userID string) ([]model.Rating, error)
	UpdateRating(rating *model.Rating) error
	DeleteRating(id string) error
	AddRatingResponse(ratingID string, response model.RatingResponse) error
}

type ratingService struct {
	repo repository.RatingRepository
}

func NewRatingService(repo repository.RatingRepository) RatingService {
	return &ratingService{repo: repo}
}

func (s *ratingService) CreateRating(rating *model.Rating) error {
	// Generate ID and timestamps
	rating.ID = uuid.New().String()
	rating.CreatedAt = time.Now()
	rating.UpdatedAt = time.Now()

	return s.repo.Create(rating)
}

func (s *ratingService) GetRatingByID(id string) (*model.Rating, error) {
	return s.repo.FindByID(id)
}

func (s *ratingService) GetAllRatings() ([]model.Rating, error) {
	return s.repo.FindAll()
}

func (s *ratingService) GetRatingsByProductID(productID string, page, limit int) ([]model.Rating, int64, error) {
	skip := (page - 1) * limit
	return s.repo.FindByProductID(productID, skip, limit)
}

func (s *ratingService) GetRatingsByUserID(userID string) ([]model.Rating, error) {
	return s.repo.FindByUserID(userID)
}

func (s *ratingService) UpdateRating(rating *model.Rating) error {
	return s.repo.Update(rating)
}

func (s *ratingService) DeleteRating(id string) error {
	return s.repo.Delete(id)
}

func (s *ratingService) AddRatingResponse(ratingID string, response model.RatingResponse) error {
	// Generate ID for the response
	response.ID = uuid.New().String()
	return s.repo.AddRatingResponse(ratingID, response)
}
