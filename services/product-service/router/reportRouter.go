package router

import (
	"product-service/controller"
	"product-service/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterReportRoutes(rg *gin.RouterGroup, ctrl *controller.ReportController) {
	reportGroup := rg.Group("")
	{
		// Customer routes - require customer authentication
		reportGroup.POST("/report", middleware.RequireCustomer(), ctrl.CreateReport)
		reportGroup.GET("/report/my", middleware.RequireCustomer(), ctrl.GetMyReports)

		// Seller/Admin routes - require seller authentication
		reportGroup.GET("/report/product/:productId", middleware.RequireSellerOrAdmin(), ctrl.GetReportsByProductID)
		reportGroup.GET("/report/:id", middleware.RequireSellerOrAdmin(), ctrl.GetReportByID)
		reportGroup.PUT("/report/:id/status", middleware.RequireAdmin(), ctrl.UpdateReportStatus)
		reportGroup.DELETE("/report/:id", middleware.RequireAdmin(), ctrl.DeleteReport)

		// Admin only routes
		reportGroup.GET("/report", middleware.RequireAdmin(), ctrl.GetAllReports)
	}
}
