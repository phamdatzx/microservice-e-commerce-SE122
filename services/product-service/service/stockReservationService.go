package service

import (
	"fmt"
	"product-service/dto"
	"product-service/model"
	"product-service/repository"
)

type StockReservationService interface {
	ReserveStock(orderID string, items []dto.ReserveStockItem) error
	ReleaseStock(orderID string) error
}

type stockReservationService struct {
	stockReservationRepo repository.StockReservationRepository
	productRepo          repository.ProductRepository
}

func NewStockReservationService(
	stockReservationRepo repository.StockReservationRepository,
	productRepo repository.ProductRepository,
) StockReservationService {
	return &stockReservationService{
		stockReservationRepo: stockReservationRepo,
		productRepo:          productRepo,
	}
}

func (s *stockReservationService) ReserveStock(orderID string, items []dto.ReserveStockItem) error {
	// 1. Validate and fetch variant details
	variantIDs := make([]string, len(items))
	for i, item := range items {
		variantIDs[i] = item.VariantID
	}

	variantToProduct, err := s.productRepo.FindVariantsByIds(variantIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch variants: %w", err)
	}

	// 2. Verify all variants exist and have sufficient stock
	for _, item := range items {
		product, exists := variantToProduct[item.VariantID]
		if !exists {
			return fmt.Errorf("variant %s not found", item.VariantID)
		}

		// Find the specific variant in the product
		var variant *model.Variant
		for i := range product.Variants {
			if product.Variants[i].ID == item.VariantID {
				variant = &product.Variants[i]
				break
			}
		}

		if variant == nil {
			return fmt.Errorf("variant %s not found in product", item.VariantID)
		}

		if variant.Stock < item.Quantity {
			return fmt.Errorf("insufficient stock for variant %s: requested %d, available %d",
				item.VariantID, item.Quantity, variant.Stock)
		}
	}

	// 3. Create stock reservations
	reservations := make([]model.StockReservation, len(items))
	for i, item := range items {
		reservations[i] = model.StockReservation{
			OrderID:   orderID,
			VariantID: item.VariantID,
			Quantity:  item.Quantity,
			Status:    "RESERVED",
		}
	}

	if err := s.stockReservationRepo.CreateMany(reservations); err != nil {
		return fmt.Errorf("failed to create stock reservations: %w", err)
	}

	// 4. Update stock for each variant (decrease stock)
	for _, item := range items {
		product := variantToProduct[item.VariantID]
		stockDelta := -item.Quantity // Negative to decrease

		if err := s.productRepo.UpdateVariantStock(product.ID, item.VariantID, stockDelta); err != nil {
			// Rollback: try to delete the reservations we just created
			// Note: This is a simple rollback. In production, consider using MongoDB transactions
			_ = s.stockReservationRepo.UpdateStatusByOrderID(orderID, "CANCELLED")
			return fmt.Errorf("failed to update stock for variant %s: %w", item.VariantID, err)
		}
	}

	return nil
}

func (s *stockReservationService) ReleaseStock(orderID string) error {
	// 1. Find all reservations for the order with status "RESERVED"
	reservations, err := s.stockReservationRepo.FindByOrderID(orderID)
	if err != nil {
		return fmt.Errorf("failed to find reservations for order %s: %w", orderID, err)
	}

	if len(reservations) == 0 {
		return fmt.Errorf("no reservations found for order %s", orderID)
	}

	// Filter only RESERVED reservations
	var reservedItems []model.StockReservation
	for _, reservation := range reservations {
		if reservation.Status == "RESERVED" {
			reservedItems = append(reservedItems, reservation)
		}
	}

	if len(reservedItems) == 0 {
		return fmt.Errorf("no reserved stock found for order %s", orderID)
	}

	// 2. Get variant details to find product IDs
	variantIDs := make([]string, len(reservedItems))
	for i, reservation := range reservedItems {
		variantIDs[i] = reservation.VariantID
	}

	variantToProduct, err := s.productRepo.FindVariantsByIds(variantIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch variants: %w", err)
	}

	// 3. Increase stock for each variant based on reservation quantity
	for _, reservation := range reservedItems {
		product, exists := variantToProduct[reservation.VariantID]
		if !exists {
			return fmt.Errorf("variant %s not found", reservation.VariantID)
		}

		stockDelta := reservation.Quantity // Positive to increase

		if err := s.productRepo.UpdateVariantStock(product.ID, reservation.VariantID, stockDelta); err != nil {
			return fmt.Errorf("failed to release stock for variant %s: %w", reservation.VariantID, err)
		}
	}

	// 4. Update reservation status to "RELEASED"
	if err := s.stockReservationRepo.UpdateStatusByOrderID(orderID, "RELEASED"); err != nil {
		return fmt.Errorf("failed to update reservation status: %w", err)
	}

	return nil
}
