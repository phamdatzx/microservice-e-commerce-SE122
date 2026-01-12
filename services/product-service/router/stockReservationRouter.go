package router

import (
	"product-service/controller"

	"github.com/gin-gonic/gin"
)

func SetupStockReservationRoutes(r *gin.RouterGroup, ctrl *controller.StockReservationController) {
	stock := r.Group("/stock")
	{
		stock.POST("/reserve", ctrl.ReserveStock)
		stock.POST("/release", ctrl.ReleaseStock)
	}
}
