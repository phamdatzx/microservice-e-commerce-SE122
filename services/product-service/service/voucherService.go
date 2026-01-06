package service

import (
	"product-service/dto"
	appError "product-service/error"
	"product-service/model"
	"product-service/repository"
)

type VoucherService interface {
	CreateVoucher(sellerID string, request dto.VoucherRequest) (dto.VoucherResponse, error)
	GetVoucherByID(id string) (dto.VoucherResponse, error)
	GetVouchersBySeller(sellerID string) ([]dto.VoucherResponse, error)
	UpdateVoucher(id string, sellerID string, request dto.VoucherRequest) (dto.VoucherResponse, error)
	DeleteVoucher(id string, sellerID string) error
}

type voucherService struct {
	repo repository.VoucherRepository
}

func NewVoucherService(repo repository.VoucherRepository) VoucherService {
	return &voucherService{repo: repo}
}

func (s *voucherService) CreateVoucher(sellerID string, request dto.VoucherRequest) (dto.VoucherResponse, error) {
	voucher := &model.Voucher{
		SellerID:               sellerID,
		Code:                   request.Code,
		Name:                   request.Name,
		Description:            request.Description,
		DiscountType:           request.DiscountType,
		DiscountValue:          request.DiscountValue,
		MaxDiscountValue:       request.MaxDiscountValue,
		MinOrderValue:          request.MinOrderValue,
		ApplyScope:             request.ApplyScope,
		ApplySellerCategoryIds: request.ApplySellerCategoryIds,
		TotalQuantity:          request.TotalQuantity,
		UsageLimitPerUser:      request.UsageLimitPerUser,
		StartTime:              request.StartTime,
		EndTime:                request.EndTime,
		Status:                 request.Status,
	}

	if err := s.repo.Create(voucher); err != nil {
		return dto.VoucherResponse{}, err
	}

	return s.mapToResponse(voucher), nil
}

func (s *voucherService) GetVoucherByID(id string) (dto.VoucherResponse, error) {
	voucher, err := s.repo.FindByID(id)
	if err != nil {
		return dto.VoucherResponse{}, appError.NewAppErrorWithErr(404, "voucher not found", err)
	}
	return s.mapToResponse(voucher), nil
}

func (s *voucherService) GetVouchersBySeller(sellerID string) ([]dto.VoucherResponse, error) {
	vouchers, err := s.repo.FindBySellerID(sellerID)
	if err != nil {
		return nil, err
	}

	var responses []dto.VoucherResponse
	for _, voucher := range vouchers {
		responses = append(responses, s.mapToResponse(&voucher))
	}
	return responses, nil
}

func (s *voucherService) UpdateVoucher(id string, sellerID string, request dto.VoucherRequest) (dto.VoucherResponse, error) {
	voucher, err := s.repo.FindByID(id)
	if err != nil {
		return dto.VoucherResponse{}, appError.NewAppErrorWithErr(404, "voucher not found", err)
	}

	if voucher.SellerID != sellerID {
		return dto.VoucherResponse{}, appError.NewAppError(403, "permission denied")
	}

	voucher.Code = request.Code
	voucher.Name = request.Name
	voucher.Description = request.Description
	voucher.DiscountType = request.DiscountType
	voucher.DiscountValue = request.DiscountValue
	voucher.MaxDiscountValue = request.MaxDiscountValue
	voucher.MinOrderValue = request.MinOrderValue
	voucher.ApplyScope = request.ApplyScope
	voucher.ApplySellerCategoryIds = request.ApplySellerCategoryIds
	voucher.TotalQuantity = request.TotalQuantity
	voucher.UsageLimitPerUser = request.UsageLimitPerUser
	voucher.StartTime = request.StartTime
	voucher.EndTime = request.EndTime
	if request.Status != "" {
		voucher.Status = request.Status
	}

	if err := s.repo.Update(voucher); err != nil {
		return dto.VoucherResponse{}, err
	}

	return s.mapToResponse(voucher), nil
}

func (s *voucherService) DeleteVoucher(id string, sellerID string) error {
	voucher, err := s.repo.FindByID(id)
	if err != nil {
		return appError.NewAppErrorWithErr(404, "voucher not found", err)
	}

	if voucher.SellerID != sellerID {
		return appError.NewAppError(403, "permission denied")
	}

	return s.repo.Delete(id)
}

func (s *voucherService) mapToResponse(voucher *model.Voucher) dto.VoucherResponse {
	return dto.VoucherResponse{
		ID:                     voucher.ID,
		SellerID:               voucher.SellerID,
		Code:                   voucher.Code,
		Name:                   voucher.Name,
		Description:            voucher.Description,
		DiscountType:           voucher.DiscountType,
		DiscountValue:          voucher.DiscountValue,
		MaxDiscountValue:       voucher.MaxDiscountValue,
		MinOrderValue:          voucher.MinOrderValue,
		ApplyScope:             voucher.ApplyScope,
		ApplySellerCategoryIds: voucher.ApplySellerCategoryIds,
		TotalQuantity:          voucher.TotalQuantity,
		UsedQuantity:           voucher.UsedQuantity,
		UsageLimitPerUser:      voucher.UsageLimitPerUser,
		StartTime:              voucher.StartTime,
		EndTime:                voucher.EndTime,
		Status:                 voucher.Status,
		CreatedAt:              voucher.CreatedAt,
		UpdatedAt:              voucher.UpdatedAt,
	}
}
