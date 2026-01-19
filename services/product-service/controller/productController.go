package controller

import (
	"net/http"
	"product-service/dto"
	"product-service/error"
	"product-service/model"
	"product-service/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service service.ProductService
}

func NewProductController(service service.ProductService) *ProductController {
	return &ProductController{service: service}
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var product model.Product
	product.SellerID = ctx.GetHeader("X-User-Id")
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	// Validate product data
	if err := product.Validate(); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Validation failed", err))
		return
	}

	if err := c.service.CreateProduct(&product); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, product)
}

func (c *ProductController) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")

	// Get userID from header if available (optional - for view history tracking)
	userID := ctx.GetHeader("X-User-Id")

	product, err := c.service.GetProductByID(id, userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	products, err := c.service.GetAllProducts()
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	var product model.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}
	product.ID = id

	// Validate product data
	if err := product.Validate(); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Validation failed", err))
		return
	}

	if err := c.service.UpdateProduct(&product); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteProduct(id); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (c *ProductController) UploadProductImages(ctx *gin.Context) {
	productID := ctx.Param("id")

	// Parse multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Failed to parse form", err))
		return
	}

	// Get all files with field name "image"
	files := form.File["image"]
	if len(files) == 0 {
		ctx.Error(error.NewAppError(http.StatusBadRequest, "No images provided"))
		return
	}

	// Process upload in service layer
	productImages, err := c.service.ProcessProductImageUpload(productID, files)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Images uploaded successfully",
		"images":  productImages,
	})
}

func (c *ProductController) UploadVariantImages(ctx *gin.Context) {
	productID := ctx.Param("id")

	// Parse multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Failed to parse form", err))
		return
	}

	if len(form.File) == 0 {
		ctx.Error(error.NewAppError(http.StatusBadRequest, "No variant images provided"))
		return
	}

	// Process upload in service layer
	uploadedVariants, err := c.service.ProcessVariantImageUpload(productID, form.File)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Variant images uploaded successfully",
		"variants": uploadedVariants,
	})
}

func (c *ProductController) GetProductsBySeller(ctx *gin.Context) {
	// Get seller ID from header
	sellerID := ctx.Param("sellerId")
	if sellerID == "" {
		ctx.Error(error.NewAppError(http.StatusUnauthorized, "Seller ID not found"))
		return
	}

	// Parse query parameters
	var params dto.GetProductsQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid query parameters", err))
		return
	}

	// Set default values
	params.SetDefaults()

	// Get products from service
	response, err := c.service.GetProductsBySeller(sellerID, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ProductController) GetVariantsByIds(ctx *gin.Context) {
	// Parse request body
	var request dto.GetVariantsByIdsRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	// Get variants from service
	variants, err := c.service.GetVariantsByIds(request.VariantIDs)
	if err != nil {
		ctx.Error(err)
		return
	}

	// Build response
	response := dto.GetVariantsByIdsResponse{
		Variants: variants,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ProductController) SearchProducts(ctx *gin.Context) {
	// Parse query parameters
	var params dto.SearchProductsQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid query parameters", err))
		return
	}

	// Set default values
	params.SetDefaults()

	// Get userID from header if available (optional - for search history tracking)
	userID := ctx.GetHeader("X-User-Id")

	// Call service method
	response, err := c.service.SearchProducts(params, userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ProductController) GetRecentlyViewedProducts(ctx *gin.Context) {
	userID := ctx.GetHeader("X-User-Id")
	if userID == "" {
		ctx.Error(error.NewAppError(401, "User ID not found in header"))
		ctx.Abort()
		return
	}

	// Get limit from query parameter, default to 20
	limitStr := ctx.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	products, err := c.service.GetRecentlyViewedProducts(userID, limit)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, products)
}
