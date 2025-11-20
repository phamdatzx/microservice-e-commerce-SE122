package service

import (
	"product-service/model"
	"product-service/repository"

	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(product *model.Product) error
	GetProductByID(id uuid.UUID) (*model.Product, error)
	GetAllProducts() ([]model.Product, error)
	UpdateProduct(product *model.Product) error
	DeleteProduct(id uuid.UUID) error
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

func (s *productService) GetProductByID(id uuid.UUID) (*model.Product, error) {
	return s.repo.FindByID(id)
}

func (s *productService) GetAllProducts() ([]model.Product, error) {
	return s.repo.FindAll()
}

func (s *productService) UpdateProduct(product *model.Product) error {
	return s.repo.Update(product)
}

func (s *productService) DeleteProduct(id uuid.UUID) error {
	return s.repo.Delete(id)
}
