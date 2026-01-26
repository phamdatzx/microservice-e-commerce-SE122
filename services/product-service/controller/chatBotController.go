package controller

import (
	"net/http"
	"product-service/dto"
	appError "product-service/error"
	"product-service/service"

	"github.com/gin-gonic/gin"
)

type ChatBotController struct {
	service service.ChatBotService
}

func NewChatBotController(service service.ChatBotService) *ChatBotController {
	return &ChatBotController{service: service}
}

// AskQuestion handles POST /api/product/chatbot/ask
func (c *ChatBotController) AskQuestion(ctx *gin.Context) {
	var request dto.ChatBotRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(appError.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	answer, err := c.service.AskQuestion(request)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.ChatBotResponse{
		Answer: answer,
	})
}
