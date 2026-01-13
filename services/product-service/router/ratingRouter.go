package router

import (
	"product-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRatingRoutes(rg *gin.RouterGroup, controller controller.RatingController) {
	ratingGroup := rg.Group("/ratings")
	{
		ratingGroup.POST("", controller.CreateRating)
		ratingGroup.GET("", controller.GetAllRatings)
		ratingGroup.GET("/:id", controller.GetRatingByID)
		ratingGroup.PUT("/:id", controller.UpdateRating)
		ratingGroup.DELETE("/:id", controller.DeleteRating)
		ratingGroup.GET("/product/:productId", controller.GetRatingsByProductID)
		ratingGroup.GET("/user", controller.GetRatingsByUserID)
		ratingGroup.POST("/:id/response", controller.AddRatingResponse)
	}
}
