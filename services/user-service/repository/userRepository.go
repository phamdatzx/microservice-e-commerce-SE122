package repository

import (
	"user-service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	CheckUserExists(username string) (bool, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(id string) (*model.User, error)
	GetSellerByID(id string) (*model.User, error)
	ActivateAccount(id string) error
	Save(user *model.User) error
	UpdateUserImage(userId string, imageURL string) error
	GetSaleInfoByUserID(userId string) (*model.SaleInfo, error)
	UpdateSaleInfo(saleInfo *model.SaleInfo) error
	UpdateProductCount(userId string, increment bool) (*model.SaleInfo, error)
	UpdateFollowCount(userId string, increment bool) (*model.SaleInfo, error)
	GetAllUsers(skip, limit int, role *string, isBanned *bool) ([]model.User, int64, error)
	UpdateUserInfo(userId string, name string, phone string, email string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) CheckUserExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Addresses").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetSellerByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("SaleInfo").Preload("Addresses").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ActivateAccount(id string) error {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return err
	}

	//change active field = true
	user.IsActive = true
	r.db.Save(&user)
	return nil
}

func (r *userRepository) Save(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateUserImage(userId string, imageURL string) error {
	var user model.User
	err := r.db.First(&user, "id = ?", userId).Error
	if err != nil {
		return err
	}

	user.Image = imageURL
	return r.db.Save(&user).Error
}

func (r *userRepository) GetSaleInfoByUserID(userId string) (*model.SaleInfo, error) {
	var saleInfo model.SaleInfo
	err := r.db.First(&saleInfo, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &saleInfo, nil
}

func (r *userRepository) UpdateSaleInfo(saleInfo *model.SaleInfo) error {
	return r.db.Save(saleInfo).Error
}

func (r *userRepository) UpdateProductCount(userId string, increment bool) (*model.SaleInfo, error) {
	// Get or create SaleInfo for the user
	var saleInfo model.SaleInfo
	err := r.db.First(&saleInfo, "user_id = ?", userId).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new SaleInfo if doesn't exist
			saleInfo = model.SaleInfo{
				UserID:       uuid.MustParse(userId),
				ProductCount: 0,
			}
			if err := r.db.Create(&saleInfo).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Update product count
	if increment {
		saleInfo.ProductCount++
	} else {
		if saleInfo.ProductCount > 0 {
			saleInfo.ProductCount--
		}
	}

	// Save updated SaleInfo
	if err := r.db.Save(&saleInfo).Error; err != nil {
		return nil, err
	}

	return &saleInfo, nil
}

func (r *userRepository) UpdateFollowCount(userId string, increment bool) (*model.SaleInfo, error) {
	// Get or create SaleInfo for the user
	var saleInfo model.SaleInfo
	err := r.db.First(&saleInfo, "user_id = ?", userId).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new SaleInfo if doesn't exist
			saleInfo = model.SaleInfo{
				UserID:      uuid.MustParse(userId),
				FollowCount: 0,
			}
			if err := r.db.Create(&saleInfo).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Update follow count
	if increment {
		saleInfo.FollowCount++
	} else {
		if saleInfo.FollowCount > 0 {
			saleInfo.FollowCount--
		}
	}

	// Save updated SaleInfo
	if err := r.db.Save(&saleInfo).Error; err != nil {
		return nil, err
	}

	return &saleInfo, nil
}

func (r *userRepository) GetAllUsers(skip, limit int, role *string, isBanned *bool) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})

	// Apply role filter if provided
	if role != nil && *role != "" {
		query = query.Where("role = ?", *role)
	}

	// Apply isBanned filter if provided
	if isBanned != nil {
		query = query.Where("is_banned = ?", *isBanned)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with preloads
	if err := query.Preload("Addresses").Preload("SaleInfo").
		Offset(skip).Limit(limit).
		Order("id DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) UpdateUserInfo(userId string, name string, phone string, email string) error {
	var user model.User
	err := r.db.First(&user, "id = ?", userId).Error
	if err != nil {
		return err
	}

	if name != "" {
		user.Name = name
	}
	if phone != "" {
		user.Phone = phone
	}
	if email != "" {
		user.Email = email
	}
	return r.db.Save(&user).Error
}
