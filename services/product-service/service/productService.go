package service

import (
	"fmt"
	"mime/multipart"
	"product-service/model"
	"product-service/repository"
	"product-service/utils"
)

type ProductService interface {
	CreateProduct(product *model.Product) error
	GetProductByID(id string) (*model.Product, error)
	GetAllProducts() ([]model.Product, error)
	UpdateProduct(product *model.Product) error
	DeleteProduct(id string) error
	ProcessProductImageUpload(productID string, files []*multipart.FileHeader) ([]model.ProductImages, error)
	ProcessVariantImageUpload(productID string, fileMap map[string][]*multipart.FileHeader) (map[string]string, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(product *model.Product) error {
	return s.repo.Create(product)
}

func (s *productService) GetProductByID(id string) (*model.Product, error) {
	return s.repo.FindByID(id)
}

func (s *productService) GetAllProducts() ([]model.Product, error) {
	return s.repo.FindAll()
}

func (s *productService) UpdateProduct(product *model.Product) error {
	return s.repo.Update(product)
}

func (s *productService) DeleteProduct(id string) error {
	return s.repo.Delete(id)
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
		return nil, err
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
