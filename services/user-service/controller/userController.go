package controller

import (
	"net/http"
	"strconv"
	"strings"
	"user-service/dto"
	appError "user-service/error"
	"user-service/service"
	"user-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) Register(ctx *gin.Context) {
	var request dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ validate struct
	if err := validate.Struct(request); err != nil {
		// lấy chi tiết lỗi
		var errors string
		for _, err := range err.(validator.ValidationErrors) {
			errors += err.Field() + " is invalid: " + err.Tag() + ", "
		}

		_ = ctx.Error(appError.NewAppError(400, errors))
		ctx.Abort()
		return
	}

	resultUser, err := c.service.RegisterUser(request)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Register success", resultUser)
}

func (c *UserController) Login(ctx *gin.Context) {
	var request dto.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := c.service.Login(request)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Login successfully", response)

}

// auth-forward verify controller
func (c *UserController) Verify(ctx *gin.Context) {
	//get auth header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.Error(appError.NewAppError(401, "Authorization header not provided"))
		ctx.Abort()
		return
	}

	//remove "Bearer" prefix
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		ctx.Error(appError.NewAppError(401, "Invalid Authorization header format"))
		ctx.Abort()
		return
	}
	token := parts[1]

	claims, err := utils.VerifyToken(token)
	if err != nil {
		ctx.Error(appError.NewAppErrorWithErr(401, "Invalid token", err))
		ctx.Abort()
		return
	}

	// Trả header để Traefik forward sang backend
	ctx.Header("X-User-Id", claims.UserID)
	ctx.Header("X-Username", claims.Username)
	ctx.Header("X-User-Role", claims.Role)

	utils.SuccessResponse(ctx, 200, "Verify successfully", nil)
}

func (c *UserController) ActivateAccount(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.Error(appError.NewAppError(401, "Token is required"))
		ctx.Abort()
		return
	}

	if err := c.service.ActivateAccount(token); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Activate account successfully", nil)
}

func (c *UserController) SendResetPasswordRequest(ctx *gin.Context) {
	var request dto.SendResetPasswordRequestDto
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := c.service.SendResetPasswordRequest(request); err != nil {
		_ = ctx.Error(appError.NewAppError(400, err.Error()))
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "SendResetPasswordRequest successfully", nil)
}

func (c *UserController) ResetPassword(ctx *gin.Context) {
	//get token from param
	token := ctx.Query("token")
	if token == "" {
		ctx.Error(appError.NewAppError(401, "Token is required"))
		ctx.Abort()
		return
	}

	//get body
	var request dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.ResetPassword(token, request); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Reset account password successfully", nil)
}

func (c *UserController) TestPrivate(ctx *gin.Context) {
	utils.SuccessResponse(ctx, 200, "Reached", nil)
}

func (c *UserController) GetMyInfo(ctx *gin.Context) {
	userId := ctx.GetHeader("X-User-Id")
	if userId == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	response, err := c.service.GetMyInfo(userId)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Get my info successfully", response)
}

func (c *UserController) CheckUsernameExists(ctx *gin.Context) {
	var request dto.CheckUsernameRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate struct
	if err := validate.Struct(request); err != nil {
		var errors string
		for _, err := range err.(validator.ValidationErrors) {
			errors += err.Field() + " is invalid: " + err.Tag() + ", "
		}
		_ = ctx.Error(appError.NewAppError(400, errors))
		ctx.Abort()
		return
	}

	response, err := c.service.CheckUsernameExists(request)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Check username successfully", response)
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	userId := ctx.Param("id")
	if userId == "" {
		ctx.Error(appError.NewAppError(400, "User ID is required"))
		ctx.Abort()
		return
	}

	response, err := c.service.GetUserByID(userId)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Get user successfully", response)
}

func (c *UserController) GetSellerByID(ctx *gin.Context) {
	sellerId := ctx.Param("id")
	if sellerId == "" {
		ctx.Error(appError.NewAppError(400, "Seller ID is required"))
		ctx.Abort()
		return
	}

	// Get user ID from header (optional - for checking follow status)
	userId := ctx.GetHeader("X-User-Id")

	response, err := c.service.GetSellerByID(sellerId, userId)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Get seller successfully", response)
}

func (c *UserController) UploadUserImage(ctx *gin.Context) {
	// Get user ID from header
	userId := ctx.GetHeader("X-User-Id")
	if userId == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	// Parse multipart form
	file, fileHeader, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.Error(appError.NewAppErrorWithErr(400, "Failed to get image file", err))
		ctx.Abort()
		return
	}
	defer file.Close()

	// Validate file type (optional but recommended)
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/jpg" && contentType != "image/webp" {
		ctx.Error(appError.NewAppError(400, "Invalid file type. Only JPEG, PNG, and WebP images are allowed"))
		ctx.Abort()
		return
	}

	// Call service to upload image
	imageURL, err := c.service.UploadUserImage(userId, file, fileHeader)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Image uploaded successfully", gin.H{
		"image_url": imageURL,
	})
}

func (c *UserController) UpdateSellerRating(ctx *gin.Context) {
	var request dto.UpdateRatingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate struct
	if err := validate.Struct(request); err != nil {
		var errors string
		for _, err := range err.(validator.ValidationErrors) {
			errors += err.Field() + " is invalid: " + err.Tag() + ", "
		}
		_ = ctx.Error(appError.NewAppError(400, errors))
		ctx.Abort()
		return
	}

	// Validate old_star for update operation
	if request.Operation == "update" && request.OldStar == 0 {
		ctx.Error(appError.NewAppError(400, "old_star is required for update operation"))
		ctx.Abort()
		return
	}

	response, err := c.service.UpdateSellerRating(request)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Seller rating updated successfully", response)
}

func (c *UserController) UpdateProductCount(ctx *gin.Context) {
	var request dto.UpdateProductCountRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate struct
	if err := validate.Struct(request); err != nil {
		var errors string
		for _, err := range err.(validator.ValidationErrors) {
			errors += err.Field() + " is invalid: " + err.Tag() + ", "
		}
		_ = ctx.Error(appError.NewAppError(400, errors))
		ctx.Abort()
		return
	}

	response, err := c.service.UpdateProductCount(request)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Product count updated successfully", response)
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	// Parse pagination parameters
	page := 1
	limit := 10

	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Parse role filter (optional)
	var role *string
	if roleStr := ctx.Query("role"); roleStr != "" {
		role = &roleStr
	}

	// Parse isBanned filter (optional)
	var isBanned *bool
	if isBannedStr := ctx.Query("is_banned"); isBannedStr != "" {
		banned := isBannedStr == "true"
		isBanned = &banned
	}

	response, err := c.service.GetAllUsers(page, limit, role, isBanned)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Get all users successfully", response)
}

type SetUserBannedRequest struct {
	IsBanned bool `json:"is_banned"`
}

func (c *UserController) SetUserBanned(ctx *gin.Context) {
	userId := ctx.Param("id")
	if userId == "" {
		ctx.Error(appError.NewAppError(400, "User ID is required"))
		ctx.Abort()
		return
	}

	var request SetUserBannedRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.SetUserBanned(userId, request.IsBanned); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	message := "User unbanned successfully"
	if request.IsBanned {
		message = "User banned successfully"
	}

	utils.SuccessResponse(ctx, 200, message, nil)
}

func (c *UserController) UpdateUserInfo(ctx *gin.Context) {
	// Get user ID from header
	userId := ctx.GetHeader("X-User-Id")
	if userId == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	var request dto.UpdateUserInfoRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate struct
	if err := validate.Struct(request); err != nil {
		var errors string
		for _, err := range err.(validator.ValidationErrors) {
			errors += err.Field() + " is invalid: " + err.Tag() + ", "
		}
		_ = ctx.Error(appError.NewAppError(400, errors))
		ctx.Abort()
		return
	}

	if err := c.service.UpdateUserInfo(userId, request); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "User info updated successfully", nil)
}
