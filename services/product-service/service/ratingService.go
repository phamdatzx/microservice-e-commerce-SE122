package service

import (
	"net/http"
	"product-service/client"
	appError "product-service/error"
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
	repo        repository.RatingRepository
	orderClient *client.OrderServiceClient
	userClient  *client.UserServiceClient
}

func NewRatingService(repo repository.RatingRepository, orderClient *client.OrderServiceClient, userClient *client.UserServiceClient) RatingService {
	return &ratingService{
		repo:        repo,
		orderClient: orderClient,
		userClient:  userClient,
	}
}

func (s *ratingService) CreateRating(rating *model.Rating) error {
	// Fetch user info from user-service
	userInfo, err := s.userClient.GetUserByID(rating.User.ID)
	if err != nil {
		return appError.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to fetch user information", err)
	}

	// Populate user data
	rating.User.Name = userInfo.Name
	rating.User.Email = userInfo.Email
	rating.User.Image = userInfo.Image
	rating.User.Phone = userInfo.Phone

	// Validate if the user has purchased the variant
	hasPurchased, err := s.orderClient.VerifyVariantPurchase(rating.User.ID, rating.ProductID, rating.VariantID)
	if err != nil {
		return appError.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to verify purchase", err)
	}

	if !hasPurchased {
		return appError.NewAppError(http.StatusForbidden, "You can only rate products you have purchased")
	}

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
