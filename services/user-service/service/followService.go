package service

import (
	"user-service/dto"
	appError "user-service/error"
	"user-service/model"
	"user-service/repository"

	"github.com/google/uuid"
)

type FollowService interface {
	FollowSeller(userID string, sellerID string) error
	UnfollowSeller(userID string, sellerID string) error
	IsFollowing(userID string, sellerID string) (bool, error)
	GetFollowers(sellerID string) ([]dto.UserFollowResponse, error)
	GetFollowing(userID string) ([]dto.UserFollowResponse, error)
	GetFollowCount(sellerID string) (int64, error)
}

type followService struct {
	repo     repository.FollowRepository
	userRepo repository.UserRepository
}

func NewFollowService(repo repository.FollowRepository, userRepo repository.UserRepository) FollowService {
	return &followService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *followService) FollowSeller(userID string, sellerID string) error {
	// Validate UUIDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return appError.NewAppError(400, "invalid user id")
	}

	sellerUUID, err := uuid.Parse(sellerID)
	if err != nil {
		return appError.NewAppError(400, "invalid seller id")
	}

	// Check if user is trying to follow themselves
	if userID == sellerID {
		return appError.NewAppError(400, "cannot follow yourself")
	}

	// Check if seller exists
	seller, err := s.userRepo.GetUserByID(sellerID)
	if err != nil {
		return appError.NewAppErrorWithErr(404, "seller not found", err)
	}

	// Optional: Check if seller has seller role
	if seller.Role != "seller" {
		return appError.NewAppError(400, "user is not a seller")
	}

	// Check if already following
	isFollowing, err := s.repo.IsFollowing(userID, sellerID)
	if err != nil {
		return err
	}
	if isFollowing {
		return appError.NewAppError(400, "already following this seller")
	}

	// Create follow relationship
	follow := &model.UserFollow{
		UserID:   userUUID,
		SellerID: sellerUUID,
	}

	return s.repo.Create(follow)
}

func (s *followService) UnfollowSeller(userID string, sellerID string) error {
	// Validate UUIDs
	if _, err := uuid.Parse(userID); err != nil {
		return appError.NewAppError(400, "invalid user id")
	}

	if _, err := uuid.Parse(sellerID); err != nil {
		return appError.NewAppError(400, "invalid seller id")
	}

	// Check if following
	isFollowing, err := s.repo.IsFollowing(userID, sellerID)
	if err != nil {
		return err
	}
	if !isFollowing {
		return appError.NewAppError(400, "not following this seller")
	}

	return s.repo.Delete(userID, sellerID)
}

func (s *followService) IsFollowing(userID string, sellerID string) (bool, error) {
	// Validate UUIDs
	if _, err := uuid.Parse(userID); err != nil {
		return false, appError.NewAppError(400, "invalid user id")
	}

	if _, err := uuid.Parse(sellerID); err != nil {
		return false, appError.NewAppError(400, "invalid seller id")
	}

	return s.repo.IsFollowing(userID, sellerID)
}

func (s *followService) GetFollowers(sellerID string) ([]dto.UserFollowResponse, error) {
	// Validate UUID
	if _, err := uuid.Parse(sellerID); err != nil {
		return nil, appError.NewAppError(400, "invalid seller id")
	}

	follows, err := s.repo.GetFollowersBySellerID(sellerID)
	if err != nil {
		return nil, err
	}

	var responses []dto.UserFollowResponse
	for _, follow := range follows {
		responses = append(responses, dto.UserFollowResponse{
			ID:       follow.ID,
			UserID:   follow.User.ID,
			Username: follow.User.Username,
			FullName: follow.User.Name,
			Image:    follow.User.Image,
		})
	}

	return responses, nil
}

func (s *followService) GetFollowing(userID string) ([]dto.UserFollowResponse, error) {
	// Validate UUID
	if _, err := uuid.Parse(userID); err != nil {
		return nil, appError.NewAppError(400, "invalid user id")
	}

	follows, err := s.repo.GetFollowingByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.UserFollowResponse
	for _, follow := range follows {
		responses = append(responses, dto.UserFollowResponse{
			ID:       follow.ID,
			UserID:   follow.Seller.ID,
			Username: follow.Seller.Username,
			FullName: follow.Seller.Name,
			Image:    follow.Seller.Image,
		})
	}

	return responses, nil
}

func (s *followService) GetFollowCount(sellerID string) (int64, error) {
	// Validate UUID
	if _, err := uuid.Parse(sellerID); err != nil {
		return 0, appError.NewAppError(400, "invalid seller id")
	}

	return s.repo.GetFollowCount(sellerID)
}
