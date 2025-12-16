package service

import (
	"net/http"
	appError "product-service/error"
	"product-service/model"
	"product-service/repository"
)

type SellerCategoryService interface {
	CreateSellerCategory(sellerCategory *model.SellerCategory) error
	GetSellerCategoryByID(id string) (*model.SellerCategory, error)
	GetSellerCategoriesBySellerID(sellerID string) ([]model.SellerCategory, error)
	GetAllSellerCategories() ([]model.SellerCategory, error)
	UpdateSellerCategory(sellerCategory *model.SellerCategory) error
	DeleteSellerCategory(userId string ,id string) error
}

type sellerCategoryService struct {
	repo repository.SellerCategoryRepository
}

func NewSellerCategoryService(repo repository.SellerCategoryRepository) SellerCategoryService {
	return &sellerCategoryService{repo: repo}
}

func (s *sellerCategoryService) CreateSellerCategory(sellerCategory *model.SellerCategory) error {
	return s.repo.Create(sellerCategory)
}

func (s *sellerCategoryService) GetSellerCategoryByID(id string) (*model.SellerCategory, error) {
	return s.repo.FindByID(id)
}

func (s *sellerCategoryService) GetSellerCategoriesBySellerID(sellerID string) ([]model.SellerCategory, error) {
	return s.repo.FindBySellerID(sellerID)
}

func (s *sellerCategoryService) GetAllSellerCategories() ([]model.SellerCategory, error) {
	return s.repo.FindAll()
}

func (s *sellerCategoryService) UpdateSellerCategory(sellerCategory *model.SellerCategory) error {
	//get seller category by id
	oldSellerCategory, err := s.repo.FindByID(sellerCategory.ID)
	if err != nil {
		return err
	}

	//check if seller has update authority
	if oldSellerCategory.SellerID != sellerCategory.SellerID {
		return appError.NewAppError(http.StatusUnauthorized, "Unauthorized")
	}

	return s.repo.Update(sellerCategory)
}

func (s *sellerCategoryService) DeleteSellerCategory(userId string ,id string) error {
		//get seller category by id
	oldSellerCategory, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	//check if seller has update authority
	if oldSellerCategory.SellerID != userId {
		return appError.NewAppError(http.StatusUnauthorized, "Unauthorized")
	}

	return s.repo.Delete(id)
}
