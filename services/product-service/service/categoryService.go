package service

import (
	"product-service/model"
	"product-service/repository"
)

type CategoryService interface {
	CreateCategory(category *model.Category) error
	GetCategoryByID(id string) (*model.Category, error)
	GetAllCategories() ([]model.Category, error)
	UpdateCategory(category *model.Category) error
	DeleteCategory(id string) error
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

func (s *categoryService) UpdateCategory(category *model.Category) error {
	return s.repo.Update(category)
}

func (s *categoryService) DeleteCategory(id string) error {
	return s.repo.Delete(id)
}
