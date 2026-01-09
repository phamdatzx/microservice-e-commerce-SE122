package router

import (
	"order-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(rg *gin.RouterGroup, c controller.OrderController) {
	order := rg.Group("")
	{
		order.GET("", c.GetOrders)
		order.POST("/checkout", c.Checkout)
		order.POST("/:orderId/payment", c.CreatePayment)
		order.POST("/public/webhook/stripe", c.StripeWebhook)
		order.GET("/test", c.Test)
	}
}
