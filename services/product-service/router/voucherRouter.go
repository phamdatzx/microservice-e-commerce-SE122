package router

import (
	"product-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterVoucherRoutes(rg *gin.RouterGroup, c *controller.VoucherController) {
	voucher := rg.Group("/vouchers")
	{
		voucher.POST("", c.Create)
		voucher.GET("/:id", c.Get)
		voucher.GET("", c.List)
		voucher.PUT("/:id", c.Update)
		voucher.DELETE("/:id", c.Delete)
	}
}
