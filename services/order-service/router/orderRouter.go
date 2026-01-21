package router

import (
	"order-service/controller"
	"order-service/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(rg *gin.RouterGroup, c controller.OrderController) {
	order := rg.Group("")
	{
		order.GET("", c.GetOrders)
		order.GET("/seller",middleware.RequireSeller(), c.GetOrdersBySellerId)
		order.GET("/seller/statistics",middleware.RequireSeller(), c.GetSellerStatistics)
		order.PUT("/:orderId",middleware.RequireSeller(), c.UpdateOrderStatus)
		order.POST("/checkout", c.Checkout)
		order.POST("/instant-checkout", c.InstantCheckout)
		order.POST("/:orderId/payment", c.CreatePayment)
		order.POST("/public/webhook/stripe", c.StripeWebhook)
		order.POST("/verify-purchase", c.VerifyPurchase)
		order.GET("/test", c.Test)
	}
}

