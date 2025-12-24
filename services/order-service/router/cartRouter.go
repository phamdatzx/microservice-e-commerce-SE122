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
		cart.GET("/public/cart", middleware.RequireCustomer(), c.GetCartItems)
		cart.POST("/public/cart", middleware.RequireCustomer(), c.AddCartItem)
		cart.PUT("/public/cart/:id", middleware.RequireCustomer(), c.UpdateCartItemQuantity)
		cart.DELETE("/public/cart/:id", middleware.RequireCustomer(), c.DeleteCartItem)
	}
}
