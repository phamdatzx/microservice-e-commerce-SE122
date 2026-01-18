package service

import (
	"fmt"
	"mime/multipart"
	"product-service/client"
	"product-service/dto"
	appError "product-service/error"
	"product-service/model"
	"product-service/repository"
	"product-service/utils"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type ProductService interface {
	CreateProduct(product *model.Product) error
	GetProductByID(id string, userID string) (*model.Product, error)
	GetAllProducts() ([]model.Product, error)
	UpdateProduct(product *model.Product) error
	DeleteProduct(id string) error
	ProcessProductImageUpload(productID string, files []*multipart.FileHeader) ([]model.ProductImages, error)
	ProcessVariantImageUpload(productID string, fileMap map[string][]*multipart.FileHeader) (map[string]string, error)
	GetProductsBySeller(sellerID string, params dto.GetProductsQueryParams) (*dto.PaginatedProductsResponse, error)
	GetVariantsByIds(variantIDs []string) ([]dto.CartVariantDto, error)
	SearchProducts(params dto.SearchProductsQueryParams, userID string) (*dto.PaginatedProductsResponse, error)
}

type productService struct {
	repo              repository.ProductRepository
	userClient        *client.UserServiceClient
	searchHistoryRepo repository.SearchHistoryRepository
}

func NewProductService(repo repository.ProductRepository, userClient *client.UserServiceClient, searchHistoryRepo repository.SearchHistoryRepository) ProductService {
	return &productService{
		repo:              repo,
		userClient:        userClient,
		searchHistoryRepo: searchHistoryRepo,
	}
}

func (s *productService) CreateProduct(product *model.Product) error {
	err := s.repo.Create(product)
	if err != nil {
		return err
	}

	// Increment product count for seller
	if err := s.userClient.UpdateProductCount(product.SellerID, "increment"); err != nil {
		// Log error but don't fail the product creation
		fmt.Printf("Failed to update product count for seller %s: %v\n", product.SellerID, err)
	}

	return nil
}

func (s *productService) GetProductByID(id string, userID string) (*model.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Save view history if user is logged in
	if userID != "" {
		viewHistory := &model.ViewHistory{
			UserID:    userID,
			ProductID: id,
		}
		// Save view history asynchronously (don't block the request)
		go func() {
			if err := s.searchHistoryRepo.CreateViewHistory(viewHistory); err != nil {
				// Log error but don't fail the request
				fmt.Printf("Failed to save view history for user %s: %v\n", userID, err)
			}
		}()
	}

	return product, nil
}

func (s *productService) GetAllProducts() ([]model.Product, error) {
	return s.repo.FindAll()
}

func (s *productService) UpdateProduct(product *model.Product) error {
	return s.repo.Update(product)
}

func (s *productService) DeleteProduct(id string) error {
	// Get product to retrieve seller ID
	product, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	err = s.repo.Delete(id)
	if err != nil {
		return err
	}

	// Decrement product count for seller
	if err := s.userClient.UpdateProductCount(product.SellerID, "decrement"); err != nil {
		// Log error but don't fail the product deletion
		fmt.Printf("Failed to update product count for seller %s: %v\n", product.SellerID, err)
	}

	return nil
}

func (s *productService) ProcessProductImageUpload(productID string, files []*multipart.FileHeader) ([]model.ProductImages, error) {
	// Get current product to determine next order number
	product, err := s.repo.FindByID(productID)
	if err != nil {
		return nil, err
	}

	nextOrder := len(product.Images) + 1
	var productImages []model.ProductImages

	// Upload each image to S3 and create ProductImages
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		// Upload to S3
		imageURL, err := utils.UploadImageToS3(file, fileHeader, "products")
		if err != nil {
			return nil, fmt.Errorf("failed to upload image to S3: %w", err)
		}

		// Create ProductImages object
		productImage := model.ProductImages{
			URL:   imageURL,
			Order: nextOrder,
		}
		productImage.BeforeCreate()

		productImages = append(productImages, productImage)
		nextOrder++
	}

	// Add images to product
	if err := s.repo.AddImagesToProduct(productID, productImages); err != nil {
		return nil, fmt.Errorf("failed to add images to product: %w", err)
	}

	return productImages, nil
}

func (s *productService) ProcessVariantImageUpload(productID string, fileMap map[string][]*multipart.FileHeader) (map[string]string, error) {
	// Get product to verify variants exist
	product, err := s.repo.FindByID(productID)
	if err != nil {
		return nil, appError.NewAppErrorWithErr(404, "Error when retrieve product with id:"+productID, err)
	}

	variantUpdates := make(map[string]string)

	// Process each variant's image
	for variantID, files := range fileMap {
		// Verify variant exists in product
		variantExists := false
		for _, v := range product.Variants {
			if v.ID == variantID {
				variantExists = true
				break
			}
		}

		if !variantExists {
			return nil, fmt.Errorf("variant %s not found in product", variantID)
		}

		// Get the first file (should only be one per variant)
		if len(files) == 0 {
			continue
		}
		fileHeader := files[0]

		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file for variant %s: %w", variantID, err)
		}
		defer file.Close()

		// Upload to S3
		imageURL, err := utils.UploadImageToS3(file, fileHeader, "products/variants")
		if err != nil {
			return nil, fmt.Errorf("failed to upload image to S3 for variant %s: %w", variantID, err)
		}

		variantUpdates[variantID] = imageURL
	}

	if len(variantUpdates) == 0 {
		return nil, fmt.Errorf("no variant images provided")
	}

	// Update variant images in database
	if err := s.repo.UpdateVariantImages(productID, variantUpdates); err != nil {
		return nil, fmt.Errorf("failed to update variant images: %w", err)
	}

	return variantUpdates, nil
}

func (s *productService) GetProductsBySeller(sellerID string, params dto.GetProductsQueryParams) (*dto.PaginatedProductsResponse, error) {
	// Import bson for building filters
	filter := make(map[string]interface{})

	// Add category filter if provided
	if params.Category != "" {
		filter["category_ids"] = params.Category
	}

	// Add seller category filter if provided
	if params.SellerCategory != "" {
		filter["seller_category_ids"] = params.SellerCategory
	}

	// Add status filter if provided
	if params.Status != "" {
		filter["status"] = params.Status
	}

	// Add search filter if provided (case-insensitive regex search on name)
	if params.Search != "" {
		filter["name"] = map[string]interface{}{
			"$regex":   params.Search,
			"$options": "i",
		}
	}

	// Determine sort field and direction
	sortField := params.SortBy
	sortDirection := 1 // ascending
	if params.SortDirection == "desc" {
		sortDirection = -1
	}

	// Handle special price sorting logic
	// If sorting by price:
	// - ascending: use price.min
	// - descending: use price.max
	if sortField == "price" {
		if sortDirection == 1 { // ascending
			sortField = "price.min"
		} else { // descending
			sortField = "price.max"
		}
	}

	// Call repository method
	products, total, err := s.repo.FindBySeller(sellerID, filter, params.GetSkip(), params.Limit, sortField, sortDirection)
	if err != nil {
		return nil, err
	}

	// Build paginated response
	response := dto.NewPaginatedProductsResponse(products, total, params.Page, params.Limit)
	return response, nil
}

func (s *productService) GetVariantsByIds(variantIDs []string) ([]dto.CartVariantDto, error) {
	// Get products containing the variants
	variantToProduct, err := s.repo.FindVariantsByIds(variantIDs)
	if err != nil {
		return nil, err
	}

	// Build response DTOs
	var result []dto.CartVariantDto
	for _, variantID := range variantIDs {
		product, exists := variantToProduct[variantID]
		if !exists {
			// Variant not found, skip it
			continue
		}

		// Find the specific variant in the product
		for _, variant := range product.Variants {
			if variant.ID == variantID {
				result = append(result, dto.CartVariantDto{
					ProductName:       product.Name,
					SellerID:          product.SellerID,
					SellerCategoryIds: product.SellerCategoryIDs,
					Variant:           variant,
				})
				break
			}
		}
	}

	return result, nil
}

func (s *productService) SearchProducts(params dto.SearchProductsQueryParams, userID string) (*dto.PaginatedProductsResponse, error) {
	// Save search history if user is logged in and has search query
	if userID != "" && params.SearchQuery != "" {
		searchHistory := &model.SearchHistory{
			UserID: userID,
			Query:  params.SearchQuery,
		}
		// Save search history asynchronously (don't block the search)
		go func() {
			if err := s.searchHistoryRepo.Create(searchHistory); err != nil {
				// Log error but don't fail the search
				fmt.Printf("Failed to save search history for user %s: %v\n", userID, err)
			}
		}()
	}

	// Build filter based on query parameters
	filter := bson.M{}

	// Price range filter
	if params.MinPrice != nil || params.MaxPrice != nil {
		priceFilter := bson.M{}
		if params.MinPrice != nil {
			// Product's max price should be >= min requested price
			priceFilter["$gte"] = *params.MinPrice
		}
		if params.MaxPrice != nil {
			// Product's min price should be <= max requested price
			filter["price.min"] = bson.M{"$lte": *params.MaxPrice}
		}
		if params.MinPrice != nil {
			filter["price.max"] = priceFilter
		}
	}

	// Rating range filter
	if params.MinRating != nil || params.MaxRating != nil {
		ratingFilter := bson.M{}
		if params.MinRating != nil {
			ratingFilter["$gte"] = *params.MinRating
		}
		if params.MaxRating != nil {
			ratingFilter["$lte"] = *params.MaxRating
		}
		filter["rating"] = ratingFilter
	}

	// Category IDs filter
	if params.CategoryIDs != "" {
		// Split comma-separated category IDs
		categoryIDsArray := strings.Split(params.CategoryIDs, ",")
		// Trim whitespace from each ID
		for i := range categoryIDsArray {
			categoryIDsArray[i] = strings.TrimSpace(categoryIDsArray[i])
		}
		// Filter products that contain at least one of the category IDs
		filter["category_ids"] = bson.M{"$in": categoryIDsArray}
	}

	// Text search filter
	sortByTextScore := false
	if params.SearchQuery != "" {
		filter["$text"] = bson.M{
			"$search": params.SearchQuery,
			"$diacriticSensitive": false,
		}
		sortByTextScore = true
	}

	// Call repository method
	products, total, err := s.repo.SearchProducts(filter, params.GetSkip(), params.Limit, sortByTextScore)
	if err != nil {
		return nil, err
	}

	// Build paginated response
	response := dto.NewPaginatedProductsResponse(products, total, params.Page, params.Limit)
	return response, nil
}
