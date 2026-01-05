package controller

import (
	"net/http"
	"product-service/error"
	"product-service/model"
	"product-service/service"

	"github.com/gin-gonic/gin"
)

type SellerCategoryController struct {
	service service.SellerCategoryService
}

func NewSellerCategoryController(service service.SellerCategoryService) *SellerCategoryController {
	return &SellerCategoryController{service: service}
}

func (c *SellerCategoryController) CreateSellerCategory(ctx *gin.Context) {
	var sellerCategory model.SellerCategory
	
	// Auto-set seller_id from X-User-Id header
	sellerCategory.SellerID = ctx.GetHeader("X-User-Id")
	
	if err := ctx.ShouldBindJSON(&sellerCategory); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	// Validate seller category data
	if err := sellerCategory.Validate(); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Validation failed", err))
		return
	}

	if err := c.service.CreateSellerCategory(&sellerCategory); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, sellerCategory)
}

func (c *SellerCategoryController) GetSellerCategoryByID(ctx *gin.Context) {
	id := ctx.Param("id")

	sellerCategory, err := c.service.GetSellerCategoryByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, sellerCategory)
}

func (c *SellerCategoryController) GetAllSellerCategoriesBySellerID(ctx *gin.Context) {
	sellerID := ctx.Param("seller_id")

	sellerCategories, err := c.service.GetSellerCategoriesBySellerID(sellerID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, sellerCategories)
}

func (c *SellerCategoryController) UpdateSellerCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	sellerID := ctx.GetHeader("X-User-Id")

	var sellerCategory model.SellerCategory
	if err := ctx.ShouldBindJSON(&sellerCategory); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}
	sellerCategory.ID = id
	sellerCategory.SellerID = sellerID

	// Validate seller category data
	if err := sellerCategory.Validate(); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Validation failed", err))
		return
	}

	if err := c.service.UpdateSellerCategory(&sellerCategory); err != nil {
		if error.IsAppError(err) {
			ctx.Error(err)
			return
		}
		ctx.Error(error.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to update seller category", err))
		return
	}

	ctx.JSON(http.StatusOK, sellerCategory)
}

func (c *SellerCategoryController) DeleteSellerCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteSellerCategory(ctx.GetHeader("X-User-Id"), id); err != nil {
		if error.IsAppError(err) {
			ctx.Error(err)
			return
		}
		ctx.Error(error.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to delete seller category", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Seller category deleted successfully"})
}
