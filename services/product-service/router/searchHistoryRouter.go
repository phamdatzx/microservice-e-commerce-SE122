package router

import (
	"product-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterSearchHistoryRoutes(rg *gin.RouterGroup, c *controller.SearchHistoryController) {
	searchHistory := rg.Group("/search-history")
	{
		searchHistory.GET("", c.GetSearchHistory)
	}
}
