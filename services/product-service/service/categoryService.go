package service

import (
	"fmt"
	"mime/multipart"
	"product-service/model"
	"product-service/repository"
	"product-service/utils"
)

type CategoryService interface {
	CreateCategory(category *model.Category) error
	GetCategoryByID(id string) (*model.Category, error)
	GetAllCategories() ([]model.Category, error)
	SearchCategories(name string) ([]model.Category, error)
	UpdateCategory(category *model.Category) error
	DeleteCategory(id string) error
	ProcessCategoryCreate(name string, imageFile *multipart.FileHeader) (*model.Category, error)
	ProcessCategoryUpdate(id string, name string, imageFile *multipart.FileHeader) (*model.Category, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(category *model.Category) error {
	return s.repo.Create(category)
}

func (s *categoryService) GetCategoryByID(id string) (*model.Category, error) {
	return s.repo.FindByID(id)
}

func (s *categoryService) GetAllCategories() ([]model.Category, error) {
	return s.repo.FindAll()
}

func (s *categoryService) SearchCategories(name string) ([]model.Category, error) {
	if name == "" {
		return s.repo.FindAll()
	}
	return s.repo.FindByName(name)
}

func (s *categoryService) UpdateCategory(category *model.Category) error {
	return s.repo.Update(category)
}

func (s *categoryService) DeleteCategory(id string) error {
	return s.repo.Delete(id)
}

func (s *categoryService) ProcessCategoryCreate(name string, imageFile *multipart.FileHeader) (*model.Category, error) {
	// Validate name
	if name == "" {
		return nil, fmt.Errorf("category name is required")
	}

	category := &model.Category{
		Name: name,
	}

	// Upload image if provided
	if imageFile != nil {
		file, err := imageFile.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open image file: %w", err)
		}
		defer file.Close()

		imageURL, err := utils.UploadImageToS3(file, imageFile, "categories")
		if err != nil {
			return nil, fmt.Errorf("failed to upload image to S3: %w", err)
		}
		category.Image = imageURL
	}

	// Create category
	if err := s.repo.Create(category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (s *categoryService) ProcessCategoryUpdate(id string, name string, imageFile *multipart.FileHeader) (*model.Category, error) {
	// Get existing category
	existingCategory, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// Update name if provided
	if name != "" {
		existingCategory.Name = name
	}

	// Upload new image if provided
	if imageFile != nil {
		file, err := imageFile.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open image file: %w", err)
		}
		defer file.Close()

		imageURL, err := utils.UploadImageToS3(file, imageFile, "categories")
		if err != nil {
			return nil, fmt.Errorf("failed to upload image to S3: %w", err)
		}
		existingCategory.Image = imageURL
	}

	// Validate
	if err := existingCategory.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update category
	if err := s.repo.Update(existingCategory); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return existingCategory, nil
}
