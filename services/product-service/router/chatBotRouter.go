package router

import (
	"product-service/controller"

	"github.com/gin-gonic/gin"
)

// RegisterChatBotRoutes registers all chatbot-related routes
func RegisterChatBotRoutes(rg *gin.RouterGroup, controller *controller.ChatBotController) {
	chatbot := rg.Group("/public/chatbot")
	{
		chatbot.POST("/ask", controller.AskQuestion)
	}
}
