package service

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"product-service/client"
	appError "product-service/error"
	"product-service/model"
	"product-service/repository"
	"product-service/utils"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type RatingService interface {
	CreateRating(rating *model.Rating, files []*multipart.FileHeader) error
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
	productRepo repository.ProductRepository
	orderClient *client.OrderServiceClient
	userClient  *client.UserServiceClient
}

func NewRatingService(repo repository.RatingRepository, productRepo repository.ProductRepository, orderClient *client.OrderServiceClient, userClient *client.UserServiceClient) RatingService {
	return &ratingService{
		repo:        repo,
		productRepo: productRepo,
		orderClient: orderClient,
		userClient:  userClient,
	}
}

func (s *ratingService) CreateRating(rating *model.Rating, files []*multipart.FileHeader) error {
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

	// Not allowed if user has already rated this product
	_, err = s.repo.FindByProductIDAndUserID(rating.ProductID, rating.User.ID)
	fmt.Println("get rating by product id and user id", err)

	if err == nil {
		// Found document → user already rated
		return appError.NewAppError(http.StatusForbidden, "You can only rate products once")
	}

	if err != mongo.ErrNoDocuments {
		return err
	}

	// Upload images to S3 if provided
	if len(files) > 0 {
		var ratingImages []model.RatingImage
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				return appError.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to open image file", err)
			}
			defer file.Close()

			imageURL, err := utils.UploadImageToS3(file, fileHeader, "ratings")
			if err != nil {
				return appError.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to upload image", err)
			}
			
			ratingImages = append(ratingImages, model.RatingImage{
				ID:  uuid.New().String(),
				URL: imageURL,
			})
		}
		rating.Images = ratingImages
	}

	//update rating info in product
	go s.UpdateRatingInfoInProduct(rating.ProductID, rating)

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


func (s *ratingService) UpdateRatingInfoInProduct(productID string, rating *model.Rating) error {
	product,err := s.productRepo.FindByID(productID)
	if err != nil {
		return err
	}
	newCount := product.RateCount + 1
	product.Rating = (product.Rating*float64(product.RateCount) + float64(rating.Star)) / float64(newCount)
	product.RateCount = newCount
	return s.productRepo.Update(product)
}