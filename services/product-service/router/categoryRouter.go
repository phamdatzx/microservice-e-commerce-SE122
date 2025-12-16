package router

import (
	"product-service/controller"
	"product-service/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(rg *gin.RouterGroup, c controller.CategoryController) {
	category := rg.Group("")
	{
		// Public routes - anyone can view categories
		category.GET("/public/category/:id", c.GetCategoryByID)
		category.GET("/public/category", c.GetAllCategories)

		// Protected routes - only admin can manage global categories
		category.POST("/category", middleware.RequireAdmin(), c.CreateCategory)
		category.PUT("/category/:id", middleware.RequireAdmin(), c.UpdateCategory)
		category.DELETE("/category/:id", middleware.RequireAdmin(), c.DeleteCategory)
	}
}
