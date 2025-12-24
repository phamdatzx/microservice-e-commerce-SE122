package router

import (
	"order-service/controller"
	"order-service/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(rg *gin.RouterGroup, c controller.OrderController) {
	cart := rg.Group("")
	{
		// Protected route - require customer authentication
		cart.POST("/public/cart", middleware.RequireCustomer(), c.AddCartItem)
		cart.DELETE("/public/cart/:id", middleware.RequireCustomer(), c.DeleteCartItem)
	}
}
