package service

import (
	"fmt"
	"order-service/client"
	"order-service/dto"
	appError "order-service/error"
	"order-service/model"
	"order-service/repository"
)

type CartService interface {
	AddCartItem(userID string, request dto.AddCartItemRequest) (*dto.CartItemResponse, error)
	DeleteCartItem(userID, cartItemID string) error
	UpdateCartItemQuantity(userID, cartItemID string, quantity int) (*dto.CartItemResponse, error)
	GetCartItems(userID string) ([]dto.CartItemDetailDto, error)
}

type cartService struct {
	repo          repository.CartRepository
	productClient *client.ProductServiceClient
	userClient    *client.UserServiceClient
}

func NewCartService(
	cartRepo repository.CartRepository,
	productClient *client.ProductServiceClient,
	userClient *client.UserServiceClient,
) CartService {
	return &cartService{
		repo:          cartRepo,
		productClient: productClient,
		userClient:    userClient,
	}
}

func (s *cartService) AddCartItem(userID string, request dto.AddCartItemRequest) (*dto.CartItemResponse, error) {
	//get variant information
	productVariants, err := s.productClient.GetVariantsByIds([]string{request.VariantID})
	if err != nil {
		return nil, appError.NewAppError(500, "failed to get variant details")
	}

	if len(productVariants) == 0 {
		return nil, appError.NewAppError(404, "variant not found")
	}

	stock := productVariants[0].Variant.Stock
	if stock < request.Quantity {
		return nil, appError.NewAppError(409, "not enough stock")
	}

	
	// Check if cart item with same variant already exists
	existingItem, err := s.repo.FindCartItemByUserAndVariant(userID, request.VariantID)
	if err != nil {
		return nil, err
	}

	// If exists, update quantity
	if existingItem != nil {
		newQuantity := existingItem.Quantity + request.Quantity
		//check if newQuantity > stock
		if newQuantity > stock {
			return nil, appError.NewAppError(409, "not enough stock")
		}
		err = s.repo.UpdateCartItemQuantity(existingItem.ID, newQuantity)
		if err != nil {
			return nil, err
		}

		// Update the existing item's quantity for response
		existingItem.Quantity = newQuantity
		return dto.ToCartItemResponse(existingItem), nil
	}

	// Fetch seller info
	seller, err := s.userClient.GetUserByID(request.SellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller info: %w", err)
	}

	// If not exists, create new cart item
	newItem := &model.CartItem{
		UserID: userID,
		Seller: model.CartSeller{
			ID:       seller.ID,
			Name:     seller.Name,
			Username: seller.Username,
		},
		Product: model.CartProduct{
			ID:                request.ProductID,
			Name:              productVariants[0].ProductName,
			SellerID:          productVariants[0].SellerID,
			SellerCategoryIDs: productVariants[0].SellerCategoryIds,
		},
		Variant: model.CartVariant{
			ID:      request.VariantID,
			SKU:     productVariants[0].Variant.SKU,
			Options: productVariants[0].Variant.Options,
			Price:   productVariants[0].Variant.Price,
			Stock:   productVariants[0].Variant.Stock,
			Image:   productVariants[0].Variant.Image,
		},
		Quantity: request.Quantity,
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

func (s *cartService) DeleteCartItem(userID, cartItemID string) error {
	// Find the cart item by ID
	cartItem, err := s.repo.FindCartItemByID(cartItemID)
	if err != nil {
		return appError.NewAppError(500, "failed to get cart item")
	}

	// Check if cart item exists
	if cartItem == nil {
		return appError.NewAppError(404, "cart item not found")
	}

	// Verify that the cart item belongs to the user
	if cartItem.UserID != userID {
		return appError.NewAppError(401, "unauthorized: you can only delete your own cart items")
	}

	// Delete the cart item
	return s.repo.DeleteCartItem(cartItemID)
}

func (s *cartService) UpdateCartItemQuantity(userID, cartItemID string, quantity int) (*dto.CartItemResponse, error) {
	// Find the cart item by ID
	cartItem, err := s.repo.FindCartItemByID(cartItemID)
	if err != nil {
		return nil, err
	}

	// Check if cart item exists
	if cartItem == nil {
		return nil, appError.NewAppError(404, "cart item not found")
	}

	// Verify that the cart item belongs to the user
	if cartItem.UserID != userID {
		return nil, appError.NewAppError(401, "unauthorized: you can only update your own cart items")
	}

	// Update the quantity
	err = s.repo.UpdateCartItemQuantity(cartItemID, quantity)
	if err != nil {
		return nil, err
	}

	// Update the cart item for response
	cartItem.Quantity = quantity
	return dto.ToCartItemResponse(cartItem), nil
}

func (s *cartService) GetCartItems(userID string) ([]dto.CartItemDetailDto, error) {
	// Get cart items from repository
	cartItems, err := s.repo.FindCartItemsByUser(userID)
	if err != nil {
		return nil, err
	}

	// If no cart items, return empty array
	if len(cartItems) == 0 {
		return []dto.CartItemDetailDto{}, nil
	}

	// Extract variant IDs
	variantIDs := make([]string, len(cartItems))
	for i, item := range cartItems {
		variantIDs[i] = item.Variant.ID
	}

	// Get variant details from product-service
	productVariants, err := s.productClient.GetVariantsByIds(variantIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get variant details: %w", err)
	}

	// Create a map for quick lookup: variant_id -> ProductVariantDto
	variantMap := make(map[string]dto.ProductVariantDto)
	for _, pv := range productVariants {
		variantMap[pv.Variant.ID] = pv
	}

	// Build enriched cart item details
	result := make([]dto.CartItemDetailDto, 0, len(cartItems))
	for _, item := range cartItems {
		productVariant, exists := variantMap[item.Variant.ID]

		// Create cart item detail
		cartItemDetail := dto.CartItemDetailDto{
			ID:       item.ID,
			UserID:   item.UserID,
			Seller:   item.Seller,
			Product: dto.CartProductDto{
				ID:                item.Product.ID,
				Name:              item.Product.Name,
				SellerID:          item.Product.SellerID,
				SellerCategoryIDs: item.Product.SellerCategoryIDs,
			},
			VariantID: item.Variant.ID,
			Quantity:  item.Quantity,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}

		// Add product info if available
		if exists {
			cartItemDetail.Variant = productVariant.Variant
		}

		result = append(result, cartItemDetail)
	}

	return result, nil
}

