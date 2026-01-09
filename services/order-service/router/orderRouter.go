package router

import (
	"order-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(rg *gin.RouterGroup, c controller.OrderController) {
	order := rg.Group("")
	{
		order.POST("/checkout", c.Checkout)
		order.POST("/public/webhook/stripe", c.StripeWebhook)
	}
}
