package router

import (
	"product-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterSavedVoucherRoutes(rg *gin.RouterGroup, c *controller.SavedVoucherController) {
	savedVoucher := rg.Group("/saved-vouchers")
	{
		savedVoucher.POST("", c.SaveVoucher)
		savedVoucher.GET("", c.GetSavedVouchers)
		savedVoucher.DELETE("/:voucherId", c.UnsaveVoucher)
	}
}
