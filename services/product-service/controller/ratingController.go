package controller

import (
	"net/http"
	"product-service/error"
	"product-service/model"
	"product-service/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RatingController struct {
	service service.RatingService
}

func NewRatingController(service service.RatingService) *RatingController {
	return &RatingController{service: service}
}

func (c *RatingController) CreateRating(ctx *gin.Context) {
	var rating model.Rating
	
	// Get user ID from header
	rating.UserID = ctx.GetHeader("X-User-Id")
	if rating.UserID == "" {
		ctx.Error(error.NewAppError(http.StatusUnauthorized, "User ID not found"))
		return
	}

	if err := ctx.ShouldBindJSON(&rating); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	// Validate required fields
	if rating.ProductID == "" {
		ctx.Error(error.NewAppError(http.StatusBadRequest, "Product ID is required"))
		return
	}
	if rating.Star < 1 || rating.Star > 5 {
		ctx.Error(error.NewAppError(http.StatusBadRequest, "Star rating must be between 1 and 5"))
		return
	}

	if err := c.service.CreateRating(&rating); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, rating)
}

func (c *RatingController) GetRatingByID(ctx *gin.Context) {
	id := ctx.Param("id")

	rating, err := c.service.GetRatingByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, rating)
}

func (c *RatingController) GetAllRatings(ctx *gin.Context) {
	ratings, err := c.service.GetAllRatings()
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, ratings)
}

func (c *RatingController) GetRatingsByProductID(ctx *gin.Context) {
	productID := ctx.Param("productId")

	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	ratings, total, err := c.service.GetRatingsByProductID(productID, page, limit)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"ratings": ratings,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func (c *RatingController) GetRatingsByUserID(ctx *gin.Context) {
	// Get user ID from header
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(error.NewAppError(http.StatusUnauthorized, "User ID not found"))
		return
	}

	ratings, err := c.service.GetRatingsByUserID(userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, ratings)
}

func (c *RatingController) UpdateRating(ctx *gin.Context) {
	id := ctx.Param("id")

	var rating model.Rating
	if err := ctx.ShouldBindJSON(&rating); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}
	rating.ID = id

	// Validate star rating if provided
	if rating.Star < 1 || rating.Star > 5 {
		ctx.Error(error.NewAppError(http.StatusBadRequest, "Star rating must be between 1 and 5"))
		return
	}

	if err := c.service.UpdateRating(&rating); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, rating)
}

func (c *RatingController) DeleteRating(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteRating(id); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Rating deleted successfully"})
}

func (c *RatingController) AddRatingResponse(ctx *gin.Context) {
	ratingID := ctx.Param("id")

	var response model.RatingResponse
	if err := ctx.ShouldBindJSON(&response); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	if response.Content == "" {
		ctx.Error(error.NewAppError(http.StatusBadRequest, "Response content is required"))
		return
	}

	if err := c.service.AddRatingResponse(ratingID, response); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "Response added successfully",
		"response": response,
	})
}
