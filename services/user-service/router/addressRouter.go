package router

import (
	"user-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterAddressRoutes(rg *gin.RouterGroup, c *controller.AddressController) {
	address := rg.Group("/addresses")
	{
		address.POST("", c.Create)
		address.GET("/:id", c.Get)
		address.GET("", c.List)
		address.PUT("/:id", c.Update)
		address.DELETE("/:id", c.Delete)
	}
}
