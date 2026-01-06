package controller

import (
	"net/http"
	"user-service/dto"
	appError "user-service/error"
	"user-service/service"
	"user-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AddressController struct {
	service service.AddressService
}

func NewAddressController(service service.AddressService) *AddressController {
	return &AddressController{service: service}
}

func (c *AddressController) Create(ctx *gin.Context) {
	var request dto.AddressRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(request); err != nil {
		var errors string
		for _, err := range err.(validator.ValidationErrors) {
			errors += err.Field() + " is invalid: " + err.Tag() + ", "
		}
		_ = ctx.Error(appError.NewAppError(400, errors))
		ctx.Abort()
		return
	}

	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	response, err := c.service.CreateAddress(userID, request)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 201, "Create address successfully", response)
}

func (c *AddressController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(appError.NewAppError(400, "Address ID is required"))
		ctx.Abort()
		return
	}

	response, err := c.service.GetAddress(id)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Get address successfully", response)
}

func (c *AddressController) List(ctx *gin.Context) {
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	response, err := c.service.GetUserAddresses(userID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Get user addresses successfully", response)
}

func (c *AddressController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(appError.NewAppError(400, "Address ID is required"))
		ctx.Abort()
		return
	}

	var request dto.AddressRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(request); err != nil {
		var errors string
		for _, err := range err.(validator.ValidationErrors) {
			errors += err.Field() + " is invalid: " + err.Tag() + ", "
		}
		_ = ctx.Error(appError.NewAppError(400, errors))
		ctx.Abort()
		return
	}

	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	response, err := c.service.UpdateAddress(id, userID, request)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Update address successfully", response)
}

func (c *AddressController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(appError.NewAppError(400, "Address ID is required"))
		ctx.Abort()
		return
	}

	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	if err := c.service.DeleteAddress(id, userID); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	utils.SuccessResponse(ctx, 200, "Delete address successfully", nil)
}
