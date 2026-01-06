package router

import (
	"order-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(rg *gin.RouterGroup, c controller.OrderController) {
	order := rg.Group("/orders")
	{
		order.POST("/checkout", c.Checkout)
	}
}
