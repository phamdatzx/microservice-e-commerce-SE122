package service

import (
	"fmt"
	"order-service/dto"
	"order-service/model"
	"order-service/repository"
)

type OrderService interface {
	AddCartItem(userID string, request dto.AddCartItemRequest) (*dto.CartItemResponse, error)
	DeleteCartItem(userID, cartItemID string) error
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

func (s *orderService) DeleteCartItem(userID, cartItemID string) error {
	// Find the cart item by ID
	cartItem, err := s.repo.FindCartItemByID(cartItemID)
	if err != nil {
		return err
	}

	// Check if cart item exists
	if cartItem == nil {
		return fmt.Errorf("cart item not found")
	}

	// Verify that the cart item belongs to the user
	if cartItem.UserID != userID {
		return fmt.Errorf("unauthorized: you can only delete your own cart items")
	}

	// Delete the cart item
	return s.repo.DeleteCartItem(cartItemID)
}

