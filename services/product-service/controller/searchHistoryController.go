package controller

import (
	"net/http"
	appError "product-service/error"
	"product-service/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SearchHistoryController struct {
	service service.SearchHistoryService
}

func NewSearchHistoryController(service service.SearchHistoryService) *SearchHistoryController {
	return &SearchHistoryController{service: service}
}

func (c *SearchHistoryController) GetSearchHistory(ctx *gin.Context) {
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	// Get limit from query parameter, default to 50
	limitStr := ctx.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	searchHistory, err := c.service.GetSearchHistory(userID, limit)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, searchHistory)
}
