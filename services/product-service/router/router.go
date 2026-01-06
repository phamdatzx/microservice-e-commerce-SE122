package router

import (
	"product-service/controller"
	"product-service/middleware"

	"github.com/gin-gonic/gin"
)

// AppRouter holds all controllers for dependency injection
type AppRouter struct {
	ProductController        *controller.ProductController
	CategoryController       *controller.CategoryController
	SellerCategoryController *controller.SellerCategoryController
	VoucherController        *controller.VoucherController
}

// SetupRouter builds the main Gin router and registers all module routes
func SetupRouter(engine *gin.Engine, appRouter *AppRouter) *gin.Engine {

	// Global middlewares
	engine.Use(
		middleware.ErrorHandler(),
	)

	api := engine.Group("/api")
	{
		// Product routes
		productGroup := api.Group("/product")
		RegisterProductRoutes(productGroup, *appRouter.ProductController)

		// Category routes
		RegisterCategoryRoutes(productGroup, *appRouter.CategoryController)

		// Seller Category routes
		RegisterSellerCategoryRoutes(productGroup, *appRouter.SellerCategoryController)

		// Voucher routes
		RegisterVoucherRoutes(productGroup, appRouter.VoucherController)
	}

	return engine
}
