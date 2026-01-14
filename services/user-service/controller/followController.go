package controller

import (
	"net/http"
	appError "user-service/error"
	"user-service/service"
	"user-service/utils"

	"github.com/gin-gonic/gin"
)

type FollowController struct {
	service service.FollowService
}

func NewFollowController(service service.FollowService) *FollowController {
	return &FollowController{service: service}
}

// Follow godoc
// @Summary Follow a seller
// @Description Allows a user to follow a seller
// @Tags follow
// @Accept json
// @Produce json
// @Param seller_id path string true "Seller ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /follows/{seller_id} [post]
func (c *FollowController) Follow(ctx *gin.Context) {
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	sellerID := ctx.Param("seller_id")
	if sellerID == "" {
		ctx.Error(appError.NewAppError(400, "Seller ID is required"))
		ctx.Abort()
		return
	}

	err := c.service.FollowSeller(userID, sellerID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Follow seller successfully", nil)
}

// Unfollow godoc
// @Summary Unfollow a seller
// @Description Allows a user to unfollow a seller
// @Tags follow
// @Accept json
// @Produce json
// @Param seller_id path string true "Seller ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /follows/{seller_id} [delete]
func (c *FollowController) Unfollow(ctx *gin.Context) {
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	sellerID := ctx.Param("seller_id")
	if sellerID == "" {
		ctx.Error(appError.NewAppError(400, "Seller ID is required"))
		ctx.Abort()
		return
	}

	err := c.service.UnfollowSeller(userID, sellerID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Unfollow seller successfully", nil)
}

// CheckFollowing godoc
// @Summary Check if following a seller
// @Description Check if the current user is following a specific seller
// @Tags follow
// @Accept json
// @Produce json
// @Param seller_id path string true "Seller ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /follows/{seller_id}/check [get]
func (c *FollowController) CheckFollowing(ctx *gin.Context) {
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	sellerID := ctx.Param("seller_id")
	if sellerID == "" {
		ctx.Error(appError.NewAppError(400, "Seller ID is required"))
		ctx.Abort()
		return
	}

	isFollowing, err := c.service.IsFollowing(userID, sellerID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Check following status successfully", gin.H{
		"is_following": isFollowing,
	})
}

// GetFollowers godoc
// @Summary Get followers of a seller
// @Description Get all users who are following a specific seller
// @Tags follow
// @Accept json
// @Produce json
// @Param seller_id path string true "Seller ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /follows/{seller_id}/followers [get]
func (c *FollowController) GetFollowers(ctx *gin.Context) {
	sellerID := ctx.Param("seller_id")
	if sellerID == "" {
		ctx.Error(appError.NewAppError(400, "Seller ID is required"))
		ctx.Abort()
		return
	}

	followers, err := c.service.GetFollowers(sellerID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Get followers successfully", followers)
}

// GetFollowing godoc
// @Summary Get users that a user is following
// @Description Get all sellers that the current user is following
// @Tags follow
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /follows/following [get]
func (c *FollowController) GetFollowing(ctx *gin.Context) {
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	following, err := c.service.GetFollowing(userID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Get following list successfully", following)
}

// GetFollowCount godoc
// @Summary Get follower count of a seller
// @Description Get the total number of followers for a specific seller
// @Tags follow
// @Accept json
// @Produce json
// @Param seller_id path string true "Seller ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /follows/{seller_id}/count [get]
func (c *FollowController) GetFollowCount(ctx *gin.Context) {
	sellerID := ctx.Param("seller_id")
	if sellerID == "" {
		ctx.Error(appError.NewAppError(400, "Seller ID is required"))
		ctx.Abort()
		return
	}

	count, err := c.service.GetFollowCount(sellerID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Get follower count successfully", gin.H{
		"count": count,
	})
}
