package main

import (
	"fmt"
	"order-service/config"
	"order-service/controller"
	"order-service/repository"
	"order-service/router"
	"order-service/service"

	"github.com/gin-gonic/gin"
)

func main() {
	//

	fmt.Println("Hello World")
	// Kết nối DB
	config.ConnectDatabase()

	// Auto migrate
	//config.DB.AutoMigrate(&model.User{})

	//wiring dependencies
	// Cart dependencies
	cartRepo := repository.NewCartRepository(config.DB)
	cartService := service.NewCartService(cartRepo)
	cartController := controller.NewCartController(cartService)

	// Order dependencies
	orderRepo := repository.NewOrderRepository(config.DB)
	orderService := service.NewOrderService(orderRepo)
	orderController := controller.NewOrderController(orderService)

	r := gin.Default()
	//r.Use(cors.Default())

	// Setup routes
	router.SetupRouter(r, &router.AppRouter{
		CartController:  cartController,
		OrderController: orderController,
	})

	r.Run(":8082") // chạy server ở port 8082
}
