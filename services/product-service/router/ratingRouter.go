package router

import (
	"product-service/controller"
	"product-service/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRatingRoutes(rg *gin.RouterGroup, controller controller.RatingController) {
	ratingGroup := rg.Group("")
	{
		ratingGroup.POST("/rating", middleware.RequireCustomer(),controller.CreateRating)
		ratingGroup.GET("public/rating", controller.GetAllRatings)
		ratingGroup.GET("public/rating/:id", controller.GetRatingByID)
		ratingGroup.PUT("/rating/:id", middleware.RequireCustomer(), controller.UpdateRating)
		ratingGroup.DELETE("/rating/:id", middleware.RequireCustomer(), controller.DeleteRating)
		ratingGroup.GET("public/rating/product/:productId", controller.GetRatingsByProductID)
		ratingGroup.GET("public/rating/user", controller.GetRatingsByUserID)
		ratingGroup.POST("/rating/:id/response", middleware.RequireSeller(), controller.AddRatingResponse)
	}
}
