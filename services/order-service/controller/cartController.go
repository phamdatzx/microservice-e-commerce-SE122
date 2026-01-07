package controller

import (
	"net/http"
	"order-service/dto"
	appError "order-service/error"
	"order-service/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type CartController struct {
	service service.CartService
}

func NewCartController(service service.CartService) *CartController {
	return &CartController{service: service}
}

func (c *CartController) AddCartItem(ctx *gin.Context) {
	// Get user ID from header
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(http.StatusUnauthorized, "User ID not found"))
		return
	}

	// Parse request body
	var request dto.AddCartItemRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(appError.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	// Validate request
	if err := validate.Struct(request); err != nil {
		ctx.Error(appError.NewAppErrorWithErr(http.StatusBadRequest, "Validation failed", err))
		return
	}

	// Add cart item
	response, err := c.service.AddCartItem(userID, request)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *CartController) DeleteCartItem(ctx *gin.Context) {
	// Get user ID from header
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(http.StatusUnauthorized, "User ID not found"))
		return
	}

	// Get cart item ID from URL parameter
	cartItemID := ctx.Param("id")
	if cartItemID == "" {
		ctx.Error(appError.NewAppError(http.StatusBadRequest, "Cart item ID is required"))
		return
	}

	// Delete cart item
	err := c.service.DeleteCartItem(userID, cartItemID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cart item deleted successfully"})
}

func (c *CartController) UpdateCartItemQuantity(ctx *gin.Context) {
	// Get user ID from header
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(http.StatusUnauthorized, "User ID not found"))
		return
	}

	// Get cart item ID from URL parameter
	cartItemID := ctx.Param("id")
	if cartItemID == "" {
		ctx.Error(appError.NewAppError(http.StatusBadRequest, "Cart item ID is required"))
		return
	}

	// Parse request body
	var request dto.UpdateCartItemQuantityRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(appError.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	// Validate request
	if err := validate.Struct(request); err != nil {
		ctx.Error(appError.NewAppErrorWithErr(http.StatusBadRequest, "Validation failed", err))
		return
	}

	// Update cart item quantity
	response, err := c.service.UpdateCartItemQuantity(userID, cartItemID, request.Quantity)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *CartController) GetCartItems(ctx *gin.Context) {
	// Get user ID from header
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(http.StatusUnauthorized, "User ID not found"))
		return
	}

	// Get cart items from service
	cartItems, err := c.service.GetCartItems(userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	// Build response
	response := dto.GetCartItemsResponse{
		CartItems: cartItems,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *CartController) GetCartItemCount(ctx *gin.Context) {
	// Get user ID from header
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(http.StatusUnauthorized, "User ID not found"))
		return
	}

	// Get cart item count from service
	count, err := c.service.GetCartItemCount(userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"count": count})
}

