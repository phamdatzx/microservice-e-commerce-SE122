package service

import (
	"product-service/model"
	"product-service/repository"
)

type SearchHistoryService interface {
	GetSearchHistory(userID string, limit int) ([]model.SearchHistory, error)
}

type searchHistoryService struct {
	searchHistoryRepo repository.SearchHistoryRepository
}

func NewSearchHistoryService(searchHistoryRepo repository.SearchHistoryRepository) SearchHistoryService {
	return &searchHistoryService{
		searchHistoryRepo: searchHistoryRepo,
	}
}

func (s *searchHistoryService) GetSearchHistory(userID string, limit int) ([]model.SearchHistory, error) {
	// Default limit to 50 if not specified or if too large
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	
	return s.searchHistoryRepo.FindByUserID(userID, limit)
}
