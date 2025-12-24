package service

import (
	"order-service/dto"
	"order-service/model"
	"order-service/repository"
)

type OrderService interface {
	AddCartItem(userID string, request dto.AddCartItemRequest) (*dto.CartItemResponse, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{repo: orderRepo}
}

func (s *orderService) AddCartItem(userID string, request dto.AddCartItemRequest) (*dto.CartItemResponse, error) {
	// Check if cart item with same variant already exists
	existingItem, err := s.repo.FindCartItemByUserAndVariant(userID, request.VariantID)
	if err != nil {
		return nil, err
	}

	// If exists, update quantity
	if existingItem != nil {
		newQuantity := existingItem.Quantity + request.Quantity
		err = s.repo.UpdateCartItemQuantity(existingItem.ID, newQuantity)
		if err != nil {
			return nil, err
		}

		// Update the existing item's quantity for response
		existingItem.Quantity = newQuantity
		return dto.ToCartItemResponse(existingItem), nil
	}

	// If not exists, create new cart item
	newItem := &model.CartItem{
		UserID:    userID,
		SellerID:  request.SellerID,
		ProductID: request.ProductID,
		VariantID: request.VariantID,
		Quantity:  request.Quantity,
	}

	// Validate the cart item
	if err := newItem.Validate(); err != nil {
		return nil, err
	}

	// Create the cart item
	err = s.repo.CreateCartItem(newItem)
	if err != nil {
		return nil, err
	}

	return dto.ToCartItemResponse(newItem), nil
}
