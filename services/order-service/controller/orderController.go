package controller

import (
	"fmt"
	"order-service/dto"
	appError "order-service/error"
	"order-service/service"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	service service.OrderService
}

func NewOrderController(service service.OrderService) *OrderController {
	return &OrderController{service: service}
}

func (c *OrderController) Checkout(ctx *gin.Context) {
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		return
	}

	var request dto.CheckoutRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(appError.NewAppErrorWithErr(400, "Invalid request body", err))
		return
	}

	response, err := c.service.Checkout(userID, request)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, response)
}

func (c *OrderController) StripeWebhook(ctx *gin.Context) {

	fmt.Println("Stripe Webhook")
	fmt.Println("body:", ctx.Request.Body)
	fmt.Println("header:", ctx.Request.Header)
	ctx.JSON(200, gin.H{"message": "Stripe Webhook"})
}