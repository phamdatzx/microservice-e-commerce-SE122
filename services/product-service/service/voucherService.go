package service

import (
	"product-service/dto"
	appError "product-service/error"
	"product-service/model"
	"product-service/repository"
	"time"
)

type VoucherService interface {
	CreateVoucher(sellerID string, request dto.VoucherRequest) (dto.VoucherResponse, error)
	GetVoucherByID(id string, userID string) (dto.VoucherResponse, error)
	GetVouchersBySeller(sellerID string, userID string) ([]dto.VoucherResponse, error)
	UpdateVoucher(id string, sellerID string, request dto.VoucherRequest) (dto.VoucherResponse, error)
	DeleteVoucher(id string, sellerID string) error
	UseVoucher(userID string, voucherID string) (dto.UseVoucherResponse, error)
}

type voucherService struct {
	repo             repository.VoucherRepository
	savedVoucherRepo repository.SavedVoucherRepository
	usageRepo        repository.VoucherUsageRepository
}

func NewVoucherService(repo repository.VoucherRepository, savedVoucherRepo repository.SavedVoucherRepository, usageRepo repository.VoucherUsageRepository) VoucherService {
	return &voucherService{
		repo:             repo,
		savedVoucherRepo: savedVoucherRepo,
		usageRepo:        usageRepo,
	}
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

	return s.mapToResponse(voucher, ""), nil
}

func (s *voucherService) GetVoucherByID(id string, userID string) (dto.VoucherResponse, error) {
	voucher, err := s.repo.FindByID(id)
	if err != nil {
		return dto.VoucherResponse{}, appError.NewAppErrorWithErr(404, "voucher not found", err)
	}
	return s.mapToResponse(voucher, userID), nil
}

func (s *voucherService) GetVouchersBySeller(sellerID string, userID string) ([]dto.VoucherResponse, error) {
	vouchers, err := s.repo.FindBySellerID(sellerID)
	if err != nil {
		return nil, err
	}

	var responses []dto.VoucherResponse
	for _, voucher := range vouchers {
		responses = append(responses, s.mapToResponse(&voucher, userID))
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

	return s.mapToResponse(voucher, ""), nil
}

func (s *voucherService) DeleteVoucher(id string, sellerID string) error {
	voucher, err := s.repo.FindByID(id)
	if err != nil {
		return appError.NewAppErrorWithErr(404, "voucher not found", err)
	}

	if voucher.SellerID != sellerID {
		return appError.NewAppError(403, "permission denied")
	}

	s.savedVoucherRepo.DeleteByVoucherID(id)

	return s.repo.Delete(id)
}

func (s *voucherService) mapToResponse(voucher *model.Voucher, userID string) dto.VoucherResponse {
	response := dto.VoucherResponse{
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

	// If userID is provided, check if user has saved this voucher
	if userID != "" && s.savedVoucherRepo != nil {
		savedVoucher, err := s.savedVoucherRepo.FindByUserAndVoucher(userID, voucher.ID)
		if err == nil && savedVoucher.ID != "" {
			response.SavedVoucher = &dto.SavedVoucherInfo{
				ID:             savedVoucher.ID,
				SavedAt:        savedVoucher.SavedAt,
				UsedCount:      savedVoucher.UsedCount,
				MaxUsesAllowed: savedVoucher.MaxUsesAllowed,
				IsDeleted:      savedVoucher.IsDeleted,
			}
		}
	}

	return response
}

func (s *voucherService) UseVoucher(userID string, voucherID string) (dto.UseVoucherResponse, error) {
	// 1. Get voucher
	voucher, err := s.repo.FindByID(voucherID)
	if err != nil {
		return dto.UseVoucherResponse{
			Success: false,
			Message: "Voucher not found",
		}, appError.NewAppErrorWithErr(404, "voucher not found" + voucherID, err)
	}

	// 2. Check if voucher is active
	if voucher.Status != "ACTIVE" {
		return dto.UseVoucherResponse{
			Success: false,
			Message: "Voucher is not active",
		}, appError.NewAppError(400, "voucher is not active")
	}

	// 3. Check if voucher is within valid time range
	now := time.Now()
	if now.Before(voucher.StartTime) {
		return dto.UseVoucherResponse{
			Success: false,
			Message: "Voucher has not started yet",
		}, appError.NewAppError(400, "voucher has not started yet")
	}
	if now.After(voucher.EndTime) {
		return dto.UseVoucherResponse{
			Success: false,
			Message: "Voucher has expired",
		}, appError.NewAppError(400, "voucher has expired")
	}

	// 4. Check if voucher has remaining quantity
	if voucher.UsedQuantity >= voucher.TotalQuantity {
		return dto.UseVoucherResponse{
			Success: false,
			Message: "Voucher has been fully used",
		}, appError.NewAppError(400, "voucher has been fully used")
	}

	// 5. Check user's usage count for this voucher
	userUsageCount, err := s.usageRepo.CountByUserAndVoucher(userID, voucherID)
	if err != nil {
		return dto.UseVoucherResponse{
			Success: false,
			Message: "Failed to check user usage",
		}, err
	}

	if int(userUsageCount) >= voucher.UsageLimitPerUser {
		return dto.UseVoucherResponse{
			Success: false,
			Message: "User has reached the usage limit for this voucher",
		}, appError.NewAppError(400, "user has reached the usage limit for this voucher")
	}

	// 6. Create usage record
	usage := &model.VoucherUsage{
		UserID:    userID,
		VoucherID: voucherID,
	}

	if err := s.usageRepo.Create(usage); err != nil {
		return dto.UseVoucherResponse{
			Success: false,
			Message: "Failed to record voucher usage",
		}, err
	}

	savedVoucher, err := s.savedVoucherRepo.FindByUserAndVoucher(userID, voucherID)
	savedVoucher.UsedCount++
	if err := s.savedVoucherRepo.Update(savedVoucher); err != nil {
		return dto.UseVoucherResponse{
			Success: false,
			Message: "Failed to update saved voucher",
		}, err
	}

	// 7. Increment voucher used quantity
	voucher.UsedQuantity++
	if err := s.repo.Update(voucher); err != nil {
		return dto.UseVoucherResponse{
			Success: false,
			Message: "Failed to update voucher",
		}, err
	}

	return dto.UseVoucherResponse{
		Success:   true,
		Message:   "Voucher used successfully",
		UsageID:   usage.ID,
		UsedAt:    usage.UsedAt,
		VoucherID: voucherID,
	}, nil
}
