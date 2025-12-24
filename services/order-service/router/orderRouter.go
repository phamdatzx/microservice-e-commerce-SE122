package router

import (
	"order-service/controller"
	"order-service/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(rg *gin.RouterGroup, c controller.OrderController) {
	cart := rg.Group("/cart")
	{
		// Protected route - require customer authentication
		cart.POST("", middleware.RequireCustomer(), c.AddCartItem)
	}
}
