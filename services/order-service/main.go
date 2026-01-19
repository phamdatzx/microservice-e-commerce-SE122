package main

import (
	"fmt"
	"order-service/client"
	stripeclient "order-service/client/payment/stripe"
	"order-service/config"
	"order-service/controller"
	"order-service/repository"
	"order-service/router"
	"order-service/service"
	"os"

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
	// Initialize clients
	productClient := client.NewProductServiceClient()
	userClient := client.NewUserServiceClient()
	stripeConfig := config.NewStripeConfig()
	stripeClient := stripeclient.NewStripeClient(stripeConfig)
	GHNClient := client.NewGHNClient()
	notificationClient := client.NewNotificationServiceClient(os.Getenv("NOTIFICATION_SERVICE_URL"))

	// Initialize layers
	cartRepo := repository.NewCartRepository(config.DB)
	orderRepo := repository.NewOrderRepository(config.DB)

	cartService := service.NewCartService(cartRepo, productClient, userClient)
	orderService := service.NewOrderService(orderRepo, cartRepo, productClient, userClient, stripeClient, GHNClient, notificationClient)

	cartController := controller.NewCartController(cartService)
	orderController := controller.NewOrderController(orderService, stripeClient)

	r := gin.Default()
	//r.Use(cors.Default())

	// Setup routes
	router.SetupRouter(r, &router.AppRouter{
		CartController:  cartController,
		OrderController: orderController,
	})

	r.Run(":8085") 
}
