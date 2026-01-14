package main

import (
	"fmt"
	"user-service/config"
	"user-service/controller"
	"user-service/model"
	"user-service/repository"
	"user-service/router"
	"user-service/service"

	"github.com/gin-gonic/gin"
)

func main() {
	//

	fmt.Println("Hello World")
	// Kết nối DB
	config.ConnectDatabase()

	//s3
	config.InitS3()


	// Auto migrate
	config.DB.AutoMigrate(&model.User{}, &model.Address{}, &model.UserFollow{})

	//wiring dependencies
	userRepo := repository.NewUserRepository(config.DB)
	addressRepo := repository.NewAddressRepository(config.DB)
	followRepo := repository.NewFollowRepository(config.DB)

	userService := service.NewUserService(userRepo)
	addressService := service.NewAddressService(addressRepo)
	followService := service.NewFollowService(followRepo, userRepo)

	userController := controller.NewUserController(userService)
	addressController := controller.NewAddressController(addressService)
	followController := controller.NewFollowController(followService)

	r := gin.Default()
	//r.Use(cors.Default())

	// Setup routes
	router.SetupRouter(r, &router.AppRouter{
		UserController:    userController,
		AddressController: addressController,
		FollowController:  followController,
	})

	r.Run(":8085")
}
