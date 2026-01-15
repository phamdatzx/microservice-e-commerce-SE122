package router

import (
	"user-service/controller"

	"github.com/gin-gonic/gin"
)

func RegisterFollowRoutes(rg *gin.RouterGroup, c *controller.FollowController) {
	follow := rg.Group("")
	{
		// Follow/unfollow a seller
		follow.POST("/follow/:seller_id", c.Follow)
		follow.DELETE("/follow/:seller_id", c.Unfollow)
		
		// Check if following a seller
		follow.GET("/follow/:seller_id/check", c.CheckFollowing)
		
		// Get followers of a seller
		follow.GET("/public/follow/:seller_id/followers", c.GetFollowers)
		
		// Get follower count of a seller
		follow.GET("/public/follow/:seller_id/count", c.GetFollowCount)
		
		// Get list of sellers that current user is following
		follow.GET("/follow", c.GetFollowing)
	}
}
