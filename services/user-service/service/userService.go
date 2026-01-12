package service

import (
	"user-service/dto"
	customError "user-service/error"
	"user-service/model"
	"user-service/repository"
	"user-service/utils"
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
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{repo: userRepo}
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

