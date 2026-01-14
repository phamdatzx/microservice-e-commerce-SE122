package repository

import (
	"user-service/model"

	"gorm.io/gorm"
)

type FollowRepository interface {
	Create(follow *model.UserFollow) error
	Delete(userID, sellerID string) error
	IsFollowing(userID, sellerID string) (bool, error)
	GetFollowersBySellerID(sellerID string) ([]model.UserFollow, error)
	GetFollowingByUserID(userID string) ([]model.UserFollow, error)
	GetFollowCount(sellerID string) (int64, error)
}

type followRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Create(follow *model.UserFollow) error {
	return r.db.Create(follow).Error
}

func (r *followRepository) Delete(userID, sellerID string) error {
	return r.db.Where("user_id = ? AND seller_id = ?", userID, sellerID).
		Delete(&model.UserFollow{}).Error
}

func (r *followRepository) IsFollowing(userID, sellerID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.UserFollow{}).
		Where("user_id = ? AND seller_id = ?", userID, sellerID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *followRepository) GetFollowersBySellerID(sellerID string) ([]model.UserFollow, error) {
	var follows []model.UserFollow
	err := r.db.Preload("User").Where("seller_id = ?", sellerID).Find(&follows).Error
	if err != nil {
		return nil, err
	}
	return follows, nil
}

func (r *followRepository) GetFollowingByUserID(userID string) ([]model.UserFollow, error) {
	var follows []model.UserFollow
	err := r.db.Preload("Seller").Where("user_id = ?", userID).Find(&follows).Error
	if err != nil {
		return nil, err
	}
	return follows, nil
}

func (r *followRepository) GetFollowCount(sellerID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.UserFollow{}).Where("seller_id = ?", sellerID).Count(&count).Error
	return count, err
}
