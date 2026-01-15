package service

import (
	"fmt"
	"mime/multipart"
	"user-service/dto"
	customError "user-service/error"
	"user-service/model"
	"user-service/repository"
	"user-service/utils"

	"github.com/google/uuid"
)

type UserService interface {
	RegisterUser(user dto.RegisterRequest) (dto.UserResponse, error)
	Login(request dto.LoginRequest) (dto.LoginResponse, error)
	ActivateAccount(token string) error
	SendResetPasswordRequest(request dto.SendResetPasswordRequestDto) error
	ResetPassword(token string, request dto.ResetPasswordRequest) error
	GetMyInfo(userId string) (dto.MyInfoResponse, error)
	CheckUsernameExists(request dto.CheckUsernameRequest) (dto.CheckUsernameResponse, error)
	GetUserByID(userId string) (dto.UserResponse, error)
	GetSellerByID(sellerId string, userId string) (dto.SellerResponse, error)
	UploadUserImage(userId string, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	UpdateSellerRating(request dto.UpdateRatingRequest) (dto.UpdateRatingResponse, error)
}

type userService struct {
	repo       repository.UserRepository
	followRepo repository.FollowRepository
}

func NewUserService(userRepo repository.UserRepository, followRepo repository.FollowRepository) UserService {
	return &userService{
		repo:       userRepo,
		followRepo: followRepo,
	}
}

func (s *userService) RegisterUser(request dto.RegisterRequest) (dto.UserResponse, error) {
	exist, err := s.repo.CheckUserExists(request.Username)
	if err != nil {
		return dto.UserResponse{}, err
	}
	if exist {
		return dto.UserResponse{}, customError.NewAppError(409, "username has already existed")
	}

	//create user model
	user := model.NewUser(request)

	//hash password and save
	user.Password, _ = utils.HashPassword(user.Password)
	err = s.repo.CreateUser(user)

	//send active email
	sendMailErr := s.sendActiveEmail(*user)
	if sendMailErr != nil {
		return dto.UserResponse{}, sendMailErr
	}

	//map to dto
	resultDto := dto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Name:     user.Name,
	}

	return resultDto, err
}

func (s *userService) Login(request dto.LoginRequest) (dto.LoginResponse, error) {
	//get user model
	user, err := s.repo.GetUserByUsername(request.Username)
	if err != nil {
		return dto.LoginResponse{}, customError.NewAppErrorWithErr(401, "username not found", err)
	}

	//check password
	if utils.CheckPasswordHash(request.Password, user.Password) { //correct info
		token, err := utils.GenerateToken(user.ID.String(), user.Username, user.Role)
		if err != nil {
			return dto.LoginResponse{}, err
		}
		return dto.LoginResponse{AccessToken: token, Role: user.Role, UserId: user.ID.String()}, nil
	} else { //incorrect info
		return dto.LoginResponse{}, customError.NewAppError(401, "incorrect password")
	}
}

func (s *userService) sendActiveEmail(user model.User) error {
	token, err := utils.GenerateToken(user.ID.String(), user.Username, user.Role)
	if err != nil {
		return err
	}

	return utils.SendEmail([]string{user.Email}, "Active your account", utils.BuildActivationEmailContent(token, user.Email))
}

func (s *userService) ActivateAccount(token string) error {
	//verify and get claims
	claims, err := utils.VerifyToken(token)
	if err != nil {
		return err
	}

	return s.repo.ActivateAccount(claims.UserID)
}

func (s *userService) SendResetPasswordRequest(request dto.SendResetPasswordRequestDto) error {

	user, err := s.repo.GetUserByUsername(request.Username)
	if err != nil {
		return err
	}

	token, err := utils.GenerateToken(user.ID.String(), user.Username, user.Role)
	if err != nil {
		return err
	}

	return utils.SendEmail([]string{user.Email}, "Reset your account password", utils.BuildResetPasswordEmailContent(token, user.Email))
}

func (s *userService) ResetPassword(token string, request dto.ResetPasswordRequest) error {
	//verify and get claims
	claims, err := utils.VerifyToken(token)
	if err != nil {
		return err
	}
	//get user from db
	user, _ := s.repo.GetUserByUsername(claims.Username)
	//change new password
	user.Password, _ = utils.HashPassword(request.Password)
	return s.repo.Save(user)
}

func (s *userService) GetMyInfo(userId string) (dto.MyInfoResponse, error) {
	user, err := s.repo.GetUserByID(userId)
	if err != nil {
		return dto.MyInfoResponse{}, err
	}

	return dto.MyInfoResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Phone:    user.Phone,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		Image:    user.Image,
		IsActive: user.IsActive,
		IsVerify: user.IsVerify,
		IsBanned: user.IsBanned,
	}, nil
}

func (s *userService) CheckUsernameExists(request dto.CheckUsernameRequest) (dto.CheckUsernameResponse, error) {
	exists, err := s.repo.CheckUserExists(request.Username)
	if err != nil {
		return dto.CheckUsernameResponse{}, err
	}

	return dto.CheckUsernameResponse{Exists: exists}, nil
}

func (s *userService) GetUserByID(userId string) (dto.UserResponse, error) {
	user, err := s.repo.GetUserByID(userId)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// Find default address
	var defaultAddress *dto.AddressResponse
	for _, addr := range user.Addresses {
		if addr.Default {
			defaultAddress = &dto.AddressResponse{
				ID:           addr.ID,
				UserID:       addr.UserID,
				FullName:     addr.FullName,
				Phone:        addr.Phone,
				AddressLine:  addr.AddressLine,
				Ward:         addr.Ward,
				District:     addr.District,
				Province:     addr.Province,
				WardCode:     addr.WardCode,
				ProvinceCode: addr.ProvinceCode,
				DistrictID:   addr.DistrictID,
				ProvinceID:   addr.ProvinceID,
				Country:      addr.Country,
				Latitude:     addr.Latitude,
				Longitude:    addr.Longitude,
				Default:      addr.Default,
			}
			break
		}
	}


	return dto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Name:     user.Name,
		Phone:    user.Phone,
		Email:    user.Email,
		Image:    user.Image,
		Address:  defaultAddress,
	}, nil
}

func (s *userService) UploadUserImage(userId string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Upload image to S3
	imageURL, err := utils.UploadImageToS3(file, fileHeader, "users")
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}

	// Update user's image field in database
	err = s.repo.UpdateUserImage(userId, imageURL)
	if err != nil {
		// If database update fails, try to clean up the uploaded image
		_ = utils.DeleteImageFromS3(imageURL)
		return "", fmt.Errorf("failed to update user image: %w", err)
	}

	return imageURL, nil
}

func (s *userService) UpdateSellerRating(request dto.UpdateRatingRequest) (dto.UpdateRatingResponse, error) {
	// Get current SaleInfo or create if not exists
	saleInfo, err := s.repo.GetSaleInfoByUserID(request.SellerID)
	if err != nil {
		// If SaleInfo doesn't exist, create it
		saleInfo = &model.SaleInfo{
			UserID:        uuid.MustParse(request.SellerID),
			RatingCount:   0,
			RatingAverage: 0,
			FollowCount:   0,
		}
	}

	switch request.Operation {
	case "create":
		// Add new rating to average
		totalRating := saleInfo.RatingAverage * float64(saleInfo.RatingCount)
		saleInfo.RatingCount++
		saleInfo.RatingAverage = (totalRating + request.Star) / float64(saleInfo.RatingCount)

	case "update":
		// Replace old rating with new one
		if saleInfo.RatingCount > 0 {
			totalRating := saleInfo.RatingAverage * float64(saleInfo.RatingCount)
			totalRating = totalRating - request.OldStar + request.Star
			saleInfo.RatingAverage = totalRating / float64(saleInfo.RatingCount)
		}

	case "delete":
		// Remove rating from average
		if saleInfo.RatingCount > 0 {
			totalRating := saleInfo.RatingAverage * float64(saleInfo.RatingCount)
			saleInfo.RatingCount--
			if saleInfo.RatingCount > 0 {
				saleInfo.RatingAverage = (totalRating - request.Star) / float64(saleInfo.RatingCount)
			} else {
				saleInfo.RatingAverage = 0
			}
		}
	}

	// Save updated SaleInfo
	err = s.repo.UpdateSaleInfo(saleInfo)
	if err != nil {
		return dto.UpdateRatingResponse{}, customError.NewAppErrorWithErr(500, "Failed to update seller rating", err)
	}

	return dto.UpdateRatingResponse{
		RatingCount:   saleInfo.RatingCount,
		RatingAverage: saleInfo.RatingAverage,
	}, nil
}

func (s *userService) GetSellerByID(sellerId string, userId string) (dto.SellerResponse, error) {
	seller, err := s.repo.GetSellerByID(sellerId)
	if err != nil {
		return dto.SellerResponse{}, err
	}

	response := dto.SellerResponse{
		ID:    seller.ID.String(),
		Name:  seller.Name,
		Image: seller.Image,
	}

	// Find and include default address
	var defaultAddress *dto.AddressResponse
	for _, addr := range seller.Addresses {
		if addr.Default {
			defaultAddress = &dto.AddressResponse{
				ID:           addr.ID,
				UserID:       addr.UserID,
				FullName:     addr.FullName,
				Phone:        addr.Phone,
				AddressLine:  addr.AddressLine,
				Ward:         addr.Ward,
				District:     addr.District,
				Province:     addr.Province,
				WardCode:     addr.WardCode,
				ProvinceCode: addr.ProvinceCode,
				DistrictID:   addr.DistrictID,
				ProvinceID:   addr.ProvinceID,
				Country:      addr.Country,
				Latitude:     addr.Latitude,
				Longitude:    addr.Longitude,
				Default:      addr.Default,
			}
			break
		}
	}
	response.Address = defaultAddress

	// Include SaleInfo if exists
	if seller.SaleInfo != nil {
		isFollowing := false
		// Check if user is following this seller (only if userId is provided)
		if userId != "" {
			isFollowing, _ = s.followRepo.IsFollowing(userId, sellerId)
		}

		response.SaleInfo = &dto.SaleInfoResponse{
			FollowCount:   seller.SaleInfo.FollowCount,
			RatingCount:   seller.SaleInfo.RatingCount,
			RatingAverage: seller.SaleInfo.RatingAverage,
			IsFollowing:   isFollowing,
		}
	}

	return response, nil
}
