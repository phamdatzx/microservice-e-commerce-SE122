package service

import (
	"errors"
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
	// 1. Fetch variant → product mapping. Used for both the fast-fail pre-check
	//    and to resolve productID for each variant.
	variantIDs := make([]string, len(items))
	for i, item := range items {
		variantIDs[i] = item.VariantID
	}

	variantToProduct, err := s.productRepo.FindVariantsByIds(variantIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch variants: %w", err)
	}

	// 2. Fast-fail pre-check: ensure every variant exists and has enough stock
	//    according to the snapshot we just read. This is not the concurrency
	//    guard — that is the atomic decrement below — but it gives callers a
	//    clear error message before doing any writes.
	for _, item := range items {
		product, exists := variantToProduct[item.VariantID]
		if !exists {
			return fmt.Errorf("variant %s not found", item.VariantID)
		}

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

	// 3. Atomically decrement each variant's stock one by one.
	//    DecrementVariantStockAtomic uses a single UpdateOne with a $gte guard,
	//    so two concurrent requests competing for the same stock cannot both
	//    succeed — the second one will get ErrInsufficientStock.
	//
	//    We track how many items were successfully decremented so we can roll
	//    back exactly those items if a later decrement fails.
	for i, item := range items {
		product := variantToProduct[item.VariantID]
		if err := s.productRepo.DecrementVariantStockAtomic(product.ID, item.VariantID, item.Quantity); err != nil {
			// Roll back all decrements that already succeeded (indices 0..i-1).
			for _, done := range items[:i] {
				doneProduct := variantToProduct[done.VariantID]
				_ = s.productRepo.UpdateVariantStock(doneProduct.ID, done.VariantID, done.Quantity)
			}
			if errors.Is(err, repository.ErrInsufficientStock) {
				return fmt.Errorf("insufficient stock for variant %s", item.VariantID)
			}
			return fmt.Errorf("failed to decrement stock for variant %s: %w", item.VariantID, err)
		}
	}

	// 4. All decrements succeeded. Now persist the reservation records.
	//    If this write fails we must restore every decremented variant.
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
		for _, item := range items {
			product := variantToProduct[item.VariantID]
			_ = s.productRepo.UpdateVariantStock(product.ID, item.VariantID, item.Quantity)
		}
		return fmt.Errorf("failed to create stock reservations: %w", err)
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
