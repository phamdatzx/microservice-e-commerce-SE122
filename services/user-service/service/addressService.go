package service

import (
	"user-service/dto"
	appError "user-service/error"
	"user-service/model"
	"user-service/repository"

	"github.com/google/uuid"
)

type AddressService interface {
	CreateAddress(userID string, request dto.AddressRequest) (dto.AddressResponse, error)
	GetAddress(id string) (dto.AddressResponse, error)
	GetUserAddresses(userID string) ([]dto.AddressResponse, error)
	UpdateAddress(id string, userID string, request dto.AddressRequest) (dto.AddressResponse, error)
	DeleteAddress(id string, userID string) error
}

type addressService struct {
	repo repository.AddressRepository
}

func NewAddressService(repo repository.AddressRepository) AddressService {
	return &addressService{repo: repo}
}

func (s *addressService) CreateAddress(userID string, request dto.AddressRequest) (dto.AddressResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return dto.AddressResponse{}, appError.NewAppError(400, "invalid user id")
	}

	address := &model.Address{
		UserID:      userUUID,
		FullName:    request.FullName,
		Phone:       request.Phone,
		AddressLine: request.AddressLine,
		Ward:        request.Ward,
		District:    request.District,
		Province:    request.Province,
		Country:     request.Country,
		Latitude:    request.Latitude,
		Longitude:   request.Longitude,
	}

	if err := s.repo.Create(address); err != nil {
		return dto.AddressResponse{}, err
	}

	return s.mapToResponse(address), nil
}

func (s *addressService) GetAddress(id string) (dto.AddressResponse, error) {
	address, err := s.repo.GetByID(id)
	if err != nil {
		return dto.AddressResponse{}, appError.NewAppErrorWithErr(404, "address not found", err)
	}
	return s.mapToResponse(address), nil
}

func (s *addressService) GetUserAddresses(userID string) ([]dto.AddressResponse, error) {
	addresses, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.AddressResponse
	for _, address := range addresses {
		responses = append(responses, s.mapToResponse(&address))
	}
	return responses, nil
}

func (s *addressService) UpdateAddress(id string, userID string, request dto.AddressRequest) (dto.AddressResponse, error) {
	address, err := s.repo.GetByID(id)
	if err != nil {
		return dto.AddressResponse{}, appError.NewAppErrorWithErr(404, "address not found", err)
	}

	if address.UserID.String() != userID {
		return dto.AddressResponse{}, appError.NewAppError(403, "permission denied")
	}

	address.FullName = request.FullName
	address.Phone = request.Phone
	address.AddressLine = request.AddressLine
	address.Ward = request.Ward
	address.District = request.District
	address.Province = request.Province
	address.Country = request.Country
	address.Latitude = request.Latitude
	address.Longitude = request.Longitude

	if err := s.repo.Update(address); err != nil {
		return dto.AddressResponse{}, err
	}

	return s.mapToResponse(address), nil
}

func (s *addressService) DeleteAddress(id string, userID string) error {
	address, err := s.repo.GetByID(id)
	if err != nil {
		return appError.NewAppErrorWithErr(404, "address not found", err)
	}

	if address.UserID.String() != userID {
		return appError.NewAppError(403, "permission denied")
	}

	return s.repo.Delete(id)
}

func (s *addressService) mapToResponse(address *model.Address) dto.AddressResponse {
	return dto.AddressResponse{
		ID:          address.ID,
		UserID:      address.UserID,
		FullName:    address.FullName,
		Phone:       address.Phone,
		AddressLine: address.AddressLine,
		Ward:        address.Ward,
		District:    address.District,
		Province:    address.Province,
		Country:     address.Country,
		Latitude:    address.Latitude,
		Longitude:   address.Longitude,
	}
}
