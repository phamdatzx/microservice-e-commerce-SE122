package controller

import (
	"net/http"
	"product-service/dto"
	appError "product-service/error"
	"product-service/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type VoucherController struct {
	service service.VoucherService
}

func NewVoucherController(service service.VoucherService) *VoucherController {
	return &VoucherController{service: service}
}

func (c *VoucherController) Create(ctx *gin.Context) {
	var request dto.VoucherRequest
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

	sellerID := ctx.GetHeader("X-User-Id")
	if sellerID == "" {
		ctx.Error(appError.NewAppError(401, "Seller ID not found in header"))
		ctx.Abort()
		return
	}

	response, err := c.service.CreateVoucher(sellerID, request)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *VoucherController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(appError.NewAppError(400, "Voucher ID is required"))
		ctx.Abort()
		return
	}

	// Get userID from header (optional)
	userID := ctx.GetHeader("X-User-Id")

	response, err := c.service.GetVoucherByID(id, userID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *VoucherController) List(ctx *gin.Context) {
	sellerID := ctx.GetHeader("X-User-Id")
	if sellerID == "" {
		ctx.Error(appError.NewAppError(401, "Seller ID not found in header"))
		ctx.Abort()
		return
	}

	// userID is same as sellerID in this context
	response, err := c.service.GetVouchersBySeller(sellerID, sellerID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *VoucherController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(appError.NewAppError(400, "Voucher ID is required"))
		ctx.Abort()
		return
	}

	var request dto.VoucherRequest
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

	sellerID := ctx.GetHeader("X-User-Id")
	if sellerID == "" {
		ctx.Error(appError.NewAppError(401, "Seller ID not found in header"))
		ctx.Abort()
		return
	}

	response, err := c.service.UpdateVoucher(id, sellerID, request)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *VoucherController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(appError.NewAppError(400, "Voucher ID is required"))
		ctx.Abort()
		return
	}

	sellerID := ctx.GetHeader("X-User-Id")
	if sellerID == "" {
		ctx.Error(appError.NewAppError(401, "Seller ID not found in header"))
		ctx.Abort()
		return
	}

	if err := c.service.DeleteVoucher(id, sellerID); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}


	ctx.JSON(http.StatusOK, gin.H{"message": "Delete voucher successfully"})
}

func (c *VoucherController) GetVouchersBySellerPublic(ctx *gin.Context) {
	sellerID := ctx.Param("sellerId")
	if sellerID == "" {
		ctx.Error(appError.NewAppError(400, "Seller ID is required"))
		ctx.Abort()
		return
	}

	// For public route, userID is empty (no authentication)
	response, err := c.service.GetVouchersBySeller(sellerID, "")
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *VoucherController) UseVoucher(ctx *gin.Context) {
	var request dto.UseVoucherRequest
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

	response, err := c.service.UseVoucher(request.UserID, request.VoucherID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, response)
}
