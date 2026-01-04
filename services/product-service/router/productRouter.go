package router

import (
	"product-service/controller"
	"product-service/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(rg *gin.RouterGroup, c controller.ProductController) {
	product := rg.Group("")
	{
		// Public routes - no authentication required
		product.GET("/public/:id", c.GetProductByID)
		product.GET("/public", c.GetAllProducts)
		product.GET("/public/products/seller/:sellerId", c.GetProductsBySeller)
		product.POST("/public/variants/batch", c.GetVariantsByIds)


		// Protected routes - require seller role
		product.POST("/", middleware.RequireSeller(), c.CreateProduct)
		product.PUT("/:id", middleware.RequireSeller(), c.UpdateProduct)
		product.DELETE("/:id", middleware.RequireSeller(), c.DeleteProduct)
		product.POST("/:id/images", middleware.RequireSeller(), c.UploadProductImages)
		product.POST("/:id/variants/images", middleware.RequireSeller(), c.UploadVariantImages)
	}
}
