package router

import (
	"order-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(rg *gin.RouterGroup, c controller.OrderController) {
	_ = rg.Group("")
	{
		// TODO: Add order routes
	}
}
