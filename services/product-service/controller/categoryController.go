package controller

import (
	"net/http"
	"product-service/error"
	"product-service/model"
	"product-service/service"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	service service.CategoryService
}

func NewCategoryController(service service.CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var category model.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	// Validate category data
	if err := category.Validate(); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Validation failed", err))
		return
	}

	if err := c.service.CreateCategory(&category); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to create category", err))
		return
	}

	ctx.JSON(http.StatusCreated, category)
}

func (c *CategoryController) GetCategoryByID(ctx *gin.Context) {
	id := ctx.Param("id")

	category, err := c.service.GetCategoryByID(id)
	if err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusNotFound, "Category not found", err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (c *CategoryController) GetAllCategories(ctx *gin.Context) {
	categories, err := c.service.GetAllCategories()
	if err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to fetch categories", err))
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	var category model.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}
	category.ID = id

	// Validate category data
	if err := category.Validate(); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Validation failed", err))
		return
	}

	if err := c.service.UpdateCategory(&category); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to update category", err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteCategory(id); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to delete category", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
