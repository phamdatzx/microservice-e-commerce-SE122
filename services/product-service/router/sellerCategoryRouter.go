package router

import (
	"product-service/controller"
	"product-service/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterSellerCategoryRoutes(rg *gin.RouterGroup, c controller.SellerCategoryController) {
	sellerCategory := rg.Group("")
	{
		// Public routes
		sellerCategory.GET("/public/seller-category/:id", c.GetSellerCategoryByID)
		sellerCategory.GET("/public/seller/:seller_id/category", c.GetAllSellerCategoriesBySellerID)

		// Seller routes - manage their own categories
		sellerCategory.POST("/seller-category", middleware.RequireSeller(), c.CreateSellerCategory)
		sellerCategory.PUT("/seller-category/:id", middleware.RequireSeller(), c.UpdateSellerCategory)
		sellerCategory.DELETE("/seller-category/:id", middleware.RequireSeller(), c.DeleteSellerCategory)

	}
}
