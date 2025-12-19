package controller

import (
	"mime/multipart"
	"net/http"
	appError "product-service/error"
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
	// Parse form data
	name := ctx.PostForm("name")
	if name == "" {
		ctx.Error(appError.NewAppError(http.StatusBadRequest, "Category name is required"))
		return
	}

	// Get image file (optional)
	var imageFile *multipart.FileHeader
	file, err := ctx.FormFile("image")
	if err == nil {
		imageFile = file
	}

	// Process in service layer
	category, err := c.service.ProcessCategoryCreate(name, imageFile)
	if err != nil {
		ctx.Error(appError.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to create category", err))
		return
	}

	ctx.JSON(http.StatusCreated, category)
}

func (c *CategoryController) GetCategoryByID(ctx *gin.Context) {
	id := ctx.Param("id")

	category, err := c.service.GetCategoryByID(id)
	if err != nil {
		ctx.Error(appError.NewAppErrorWithErr(http.StatusNotFound, "Category not found", err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (c *CategoryController) GetAllCategories(ctx *gin.Context) {
	// Get optional search query parameter
	searchName := ctx.Query("name")
	
	var categories []model.Category
	var err error
	
	if searchName != "" {
		// Search by name
		categories, err = c.service.SearchCategories(searchName)
	} else {
		// Get all categories
		categories, err = c.service.GetAllCategories()
	}
	
	if err != nil {
		ctx.Error(appError.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to fetch categories", err))
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	// Parse form data
	name := ctx.PostForm("name")

	// Get image file (optional)
	var imageFile *multipart.FileHeader
	file, err := ctx.FormFile("image")
	if err == nil {
		imageFile = file
	}

	// Process in service layer
	category, err := c.service.ProcessCategoryUpdate(id, name, imageFile)
	if err != nil {
		ctx.Error(appError.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to update category", err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteCategory(id); err != nil {
		ctx.Error(appError.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to delete category", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
