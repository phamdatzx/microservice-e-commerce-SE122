package main

import (
	"fmt"
	"product-service/client"
	"product-service/config"
	"product-service/controller"
	"product-service/repository"
	"product-service/router"
	"product-service/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello Product Service")
	// Kết nối DB
	config.ConnectDatabase()

	// Initialize S3
	config.InitS3()

	// MongoDB doesn't require AutoMigrate like GORM
	// Collections will be created automatically when first document is inserted

	// Wiring dependencies - Product
	productRepo := repository.NewProductRepository(config.DB)
	userClient := client.NewUserServiceClient()
	searchHistoryRepo := repository.NewSearchHistoryRepository(config.DB)
	productService := service.NewProductService(productRepo, userClient, searchHistoryRepo)
	productController := controller.NewProductController(productService)

	// Wiring dependencies - Category
	categoryRepo := repository.NewCategoryRepository(config.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryController := controller.NewCategoryController(categoryService)

	// Wiring dependencies - SellerCategory
	sellerCategoryRepo := repository.NewSellerCategoryRepository(config.DB)
	sellerCategoryService := service.NewSellerCategoryService(sellerCategoryRepo)
	sellerCategoryController := controller.NewSellerCategoryController(sellerCategoryService)

	// Wiring dependencies - SavedVoucher (initialize repo first)
	savedVoucherRepo := repository.NewSavedVoucherRepository(config.DB)

	// Wiring dependencies - VoucherUsage
	voucherUsageRepo := repository.NewVoucherUsageRepository(config.DB)

	// Wiring dependencies - Voucher
	voucherRepo := repository.NewVoucherRepository(config.DB)
	voucherService := service.NewVoucherService(voucherRepo, savedVoucherRepo, voucherUsageRepo)
	voucherController := controller.NewVoucherController(voucherService)

	// Wiring dependencies - SavedVoucher service
	savedVoucherService := service.NewSavedVoucherService(savedVoucherRepo, voucherRepo)
	savedVoucherController := controller.NewSavedVoucherController(savedVoucherService)

	// Wiring dependencies - StockReservation
	stockReservationRepo := repository.NewStockReservationRepository(config.DB)
	stockReservationService := service.NewStockReservationService(stockReservationRepo, productRepo)
	stockReservationController := controller.NewStockReservationController(stockReservationService)

	// Wiring dependencies - Rating
	ratingRepo := repository.NewRatingRepository(config.DB)
	orderClient := client.NewOrderServiceClient()
	ratingService := service.NewRatingService(ratingRepo,productRepo, orderClient, userClient)
	ratingController := controller.NewRatingController(ratingService)

	// Wiring dependencies - SearchHistory
	searchHistoryService := service.NewSearchHistoryService(searchHistoryRepo)
	searchHistoryController := controller.NewSearchHistoryController(searchHistoryService)

	// Wiring dependencies - Report
	reportRepo := repository.NewReportRepository(config.DB)
	reportService := service.NewReportService(reportRepo, orderClient, userClient)
	reportController := controller.NewReportController(reportService)

	// Wiring dependencies - ChatBot
	chatBotClient := client.NewChatBotClient()
	chatBotService := service.NewChatBotService(chatBotClient, productRepo, categoryRepo, sellerCategoryRepo, voucherRepo, userClient)
	chatBotController := controller.NewChatBotController(chatBotService)

	r := gin.Default()
	r.Use(cors.Default())

	// Setup routes
	router.SetupRouter(r, &router.AppRouter{
		ProductController:          productController,
		CategoryController:         categoryController,
		SellerCategoryController:   sellerCategoryController,
		VoucherController:          voucherController,
		SavedVoucherController:     savedVoucherController,
		StockReservationController: stockReservationController,
		RatingController:           ratingController,
		SearchHistoryController:    searchHistoryController,
		ReportController:           reportController,
		ChatBotController:          chatBotController,
	})

	r.Run(":8085")
}
