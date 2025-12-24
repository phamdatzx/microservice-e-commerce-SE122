package router

import (
	"order-service/controller"
	"order-service/middleware"

	"github.com/gin-gonic/gin"
)

// AppRouter holds all controllers for dependency injection
type AppRouter struct {
	OrderController *controller.OrderController
}

// SetupRouter builds the main Gin router and registers all module routes
func SetupRouter(engine *gin.Engine, appRouter *AppRouter) *gin.Engine {

	// Global middlewares
	engine.Use(
		middleware.ErrorHandler(),
	)

	api := engine.Group("/api/order")
	{
		// Each controller registers its own routes
		RegisterOrderRoutes(api, *appRouter.OrderController)
	}

	//publicApi := engine.Group("/api/public")
	//{
	//	RegisterPublicUserRoutes(publicApi, *appRouter.UserController)
	//}

	return engine
}
