package router

import (
	"user-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, c controller.UserController) {
	user := rg.Group("")
	{
		user.POST("/register", c.Register)
		user.POST("/login", c.Login)
		user.GET("/verify", c.Verify)
		user.POST("/activate", c.ActivateAccount)
		user.POST("/send-reset-password", c.SendResetPasswordRequest)
		user.POST("/reset-password", c.ResetPassword)
	}
}
