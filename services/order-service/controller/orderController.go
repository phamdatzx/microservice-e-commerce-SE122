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

type OrderController struct {
	service service.OrderService
}

func NewOrderController(service service.OrderService) *OrderController {
	return &OrderController{service: service}
}

func (c *OrderController) AddCartItem(ctx *gin.Context) {
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
		ctx.Error(appError.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to add cart item", err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}
