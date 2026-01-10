package controller

import (
	"net/http"
	"product-service/dto"
	appError "product-service/error"
	"product-service/service"

	"github.com/gin-gonic/gin"
)

type SavedVoucherController struct {
	service service.SavedVoucherService
}

func NewSavedVoucherController(service service.SavedVoucherService) *SavedVoucherController {
	return &SavedVoucherController{service: service}
}

func (c *SavedVoucherController) SaveVoucher(ctx *gin.Context) {
	var request dto.SaveVoucherRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(request); err != nil {
		_ = ctx.Error(appError.NewAppError(400, "voucher_id is required"))
		ctx.Abort()
		return
	}

	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	if err := c.service.SaveVoucher(userID, request.VoucherID); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Voucher saved successfully"})
}

func (c *SavedVoucherController) GetSavedVouchers(ctx *gin.Context) {
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	response, err := c.service.GetSavedVouchers(userID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *SavedVoucherController) UnsaveVoucher(ctx *gin.Context) {
	voucherID := ctx.Param("voucherId")
	if voucherID == "" {
		ctx.Error(appError.NewAppError(400, "Voucher ID is required"))
		ctx.Abort()
		return
	}

	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(appError.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	if err := c.service.UnsaveVoucher(userID, voucherID); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Voucher unsaved successfully"})
}
