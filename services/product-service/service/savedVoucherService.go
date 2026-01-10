package service

import (
	"product-service/dto"
	appError "product-service/error"
	"product-service/model"
	"product-service/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

type SavedVoucherService interface {
	SaveVoucher(userID, voucherID string) error
	GetSavedVouchers(userID string) ([]dto.SavedVoucherResponse, error)
	UnsaveVoucher(userID, voucherID string) error
}

type savedVoucherService struct {
	savedVoucherRepo repository.SavedVoucherRepository
	voucherRepo      repository.VoucherRepository
}

func NewSavedVoucherService(savedVoucherRepo repository.SavedVoucherRepository, voucherRepo repository.VoucherRepository) SavedVoucherService {
	return &savedVoucherService{
		savedVoucherRepo: savedVoucherRepo,
		voucherRepo:      voucherRepo,
	}
}

func (s *savedVoucherService) SaveVoucher(userID, voucherID string) error {
	// Check if voucher exists and get its details
	voucher, err := s.voucherRepo.FindByID(voucherID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return appError.NewAppError(404, "voucher not found")
		}
		return err
	}

	// Check if already saved
	existingSaved, err := s.savedVoucherRepo.FindByUserAndVoucher(userID, voucherID)
	if err == nil && existingSaved.ID != "" {
		return appError.NewAppError(409, "voucher already saved")
	}

	// Save the voucher with max uses allowed from voucher's usage limit per user
	savedVoucher := &model.SavedVoucher{
		UserID:         userID,
		VoucherID:      voucherID,
		MaxUsesAllowed: voucher.UsageLimitPerUser,
	}

	return s.savedVoucherRepo.Create(savedVoucher)
}

func (s *savedVoucherService) GetSavedVouchers(userID string) ([]dto.SavedVoucherResponse, error) {
	savedVouchers, err := s.savedVoucherRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.SavedVoucherResponse
	for _, saved := range savedVouchers {
		// Get full voucher details
		voucher, err := s.voucherRepo.FindByID(saved.VoucherID)
		if err != nil {
			// Skip if voucher no longer exists
			continue
		}

		response := dto.SavedVoucherResponse{
			ID:             saved.ID,
			UserID:         saved.UserID,
			VoucherID:      saved.VoucherID,
			SavedAt:        saved.SavedAt,
			UsedCount:      saved.UsedCount,
			MaxUsesAllowed: saved.MaxUsesAllowed,
			Voucher: dto.VoucherResponse{
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
			},
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *savedVoucherService) UnsaveVoucher(userID, voucherID string) error {
	// Check if the saved voucher exists
	_, err := s.savedVoucherRepo.FindByUserAndVoucher(userID, voucherID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return appError.NewAppError(404, "saved voucher not found")
		}
		return err
	}

	return s.savedVoucherRepo.Delete(userID, voucherID)
}
