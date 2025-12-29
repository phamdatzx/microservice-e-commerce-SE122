package router

import (
	"order-service/controller"
	"order-service/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCartRoutes(rg *gin.RouterGroup, c controller.CartController) {
	cart := rg.Group("")
	{
		// Protected route - require customer authentication
		cart.GET("/cart", middleware.RequireCustomer(), c.GetCartItems)
		cart.POST("/cart", middleware.RequireCustomer(), c.AddCartItem)
		cart.PUT("/cart/:id", middleware.RequireCustomer(), c.UpdateCartItemQuantity)
		cart.DELETE("/cart/:id", middleware.RequireCustomer(), c.DeleteCartItem)
	}
}
