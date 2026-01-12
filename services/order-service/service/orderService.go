package service

import (
	"context"
	"fmt"
	"order-service/client"
	"order-service/client/payment"
	"order-service/dto"
	appError "order-service/error"
	"order-service/model"
	"order-service/repository"
	"os"
	"time"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
)

type OrderService interface {
	Checkout(userID string, request dto.CheckoutRequest) (*dto.CheckoutResponse, error)
	UpdateOrderPaymentStatus(ctx context.Context, orderID string, status string) error
	CreateCheckoutSession(order *model.Order, successURL, cancelURL string) (string, error)
	HandleCheckoutSessionCompleted(ctx context.Context, orderID, sessionID, paymentIntentID string) error
	HandlePaymentFailed(ctx context.Context, paymentIntentID string) error
	CreatePaymentForOrder(ctx context.Context, orderID string) (*dto.CreatePaymentResponse, error)
	GetUserOrders(ctx context.Context, userID string, request dto.GetOrdersRequest) (*dto.GetOrdersResponse, error)
	GetSellerOrders(ctx context.Context, sellerID string, request dto.GetOrdersBySellerRequest) (*dto.GetOrdersResponse, error)
	UpdateOrderStatus(ctx context.Context, userID string, orderID string, request dto.UpdateOrderStatusRequest) error
}

type orderService struct {
	repo          repository.OrderRepository
	cartRepo      repository.CartRepository
	productClient *client.ProductServiceClient
	userClient    *client.UserServiceClient
	paymentClient payment.PaymentClient
	GHNClient     *client.GHNClient
	clientURL     string
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	productClient *client.ProductServiceClient,
	userClient *client.UserServiceClient,
	paymentClient payment.PaymentClient,
	GHNClient *client.GHNClient,
) OrderService {
	return &orderService{
		repo:          orderRepo,
		cartRepo:      cartRepo,
		productClient: productClient,
		userClient:    userClient,
		paymentClient: paymentClient,
		GHNClient:     GHNClient,
		clientURL:     os.Getenv("CLIENT_URL"),
	}
}

func (s *orderService) Checkout(userID string, request dto.CheckoutRequest) (*dto.CheckoutResponse, error) {
	//fetch seller
	sellerFetch, err := s.userClient.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// 1. Fetch cart items
	var cartItems []*model.CartItem
	for _, itemID := range request.CartItemIDs {
		item, err := s.cartRepo.FindCartItemByID(itemID)
		if err != nil {
			return nil, err
		}
		if item == nil {
			return nil, appError.NewAppError(404, fmt.Sprintf("cart item %s not found", itemID))
		}
		if item.UserID != userID {
			return nil, appError.NewAppError(403, "cart item does not belong to user")
		}
		cartItems = append(cartItems, item)
	}

	if len(cartItems) == 0 {
		return nil, appError.NewAppError(400, "no cart items provided")
	}

	// 2. Validate all items from same seller
	sellerID := cartItems[0].Seller.ID
	for _, item := range cartItems {
		if item.Seller.ID != sellerID {
			return nil, appError.NewAppError(400, "all items must be from the same seller")
		}
	}

	// 3. Fetch User and Seller info
	user, err := s.userClient.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	seller, err := s.userClient.GetUserByID(sellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller info: %w", err)
	}

	// 4. Fetch Variant details and reverse stock
	var orderItems []model.OrderItem
	var totalAmount float64
	var variantIDs []string

	for _, item := range cartItems {
		variantIDs = append(variantIDs, item.Variant.ID)
	}

	variants, err := s.productClient.GetVariantsByIds(variantIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get variant details: %w", err)
	}

	variantMap := make(map[string]dto.ProductVariantDto)
	for _, v := range variants {
		variantMap[v.Variant.ID] = dto.ProductVariantDto{
			ProductName:       v.ProductName,
			SellerID:          v.SellerID,
			SellerCategoryIds: v.SellerCategoryIds,
			Variant: dto.VariantDto{
				ID:      v.Variant.ID,
				SKU:     v.Variant.SKU,
				Options: v.Variant.Options,
				Price:   v.Variant.Price,
				Stock:   v.Variant.Stock,
				Image:   v.Variant.Image,
			},
		}
	}

	for _, item := range cartItems {
		variantInfo, exists := variantMap[item.Variant.ID]
		if !exists {
			return nil, appError.NewAppError(404, fmt.Sprintf("variant %s not found", item.Variant.ID))
		}

		if variantInfo.Variant.Stock < item.Quantity {
			return nil, appError.NewAppError(409, fmt.Sprintf("not enough stock for product %s", variantInfo.ProductName))
		}

		itemTotal := float64(variantInfo.Variant.Price * item.Quantity)
		totalAmount += itemTotal

		// Convert options map to string
		variantName := ""
		for k, v := range variantInfo.Variant.Options {
			variantName += fmt.Sprintf("%s: %s, ", k, v)
		}
		if len(variantName) > 2 {
			variantName = variantName[:len(variantName)-2]
		}

		orderItems = append(orderItems, model.OrderItem{
			ProductID:   item.Product.ID,
			ProductName: variantInfo.ProductName,
			VariantID:   item.Variant.ID,
			VariantName: variantName,
			SKU:         variantInfo.Variant.SKU,
			Price:       variantInfo.Variant.Price,
			Quantity:    item.Quantity,
			Image:       variantInfo.Variant.Image,
		})
	}

	// Reserve stock for all items
	var reserveItems []client.ReserveStockItem
	for _, item := range cartItems {
		reserveItems = append(reserveItems, client.ReserveStockItem{
			VariantID: item.Variant.ID,
			Quantity:  item.Quantity,
		})
	}

	// Create a temporary order ID for stock reservation
	// We'll use this to track the reservation before the actual order is created
	tempOrderID := fmt.Sprintf("temp_%s_%d", userID, time.Now().UnixNano())

	if err := s.productClient.ReserveStock(tempOrderID, reserveItems); err != nil {
		return nil, fmt.Errorf("failed to reserve stock: %w", err)
	}

	// 5. Apply Voucher
	var orderVoucher *model.OrderVoucher
	if request.VoucherID != "" {
		voucher, err := s.productClient.GetVoucherByID(request.VoucherID)
		if err != nil {
			// Rollback: release reserved stock
			_ = s.productClient.ReleaseStock(tempOrderID)
			return nil, fmt.Errorf("failed to get voucher: %w", err)
		}

		// Validate voucher (basic validation, more complex logic might be needed)
		if voucher.Status != "ACTIVE" {
			// Rollback: release reserved stock
			_ = s.productClient.ReleaseStock(tempOrderID)
			return nil, appError.NewAppError(400, "voucher is not active")
		}
		if totalAmount < float64(voucher.MinOrderValue) {
			// Rollback: release reserved stock
			_ = s.productClient.ReleaseStock(tempOrderID)
			return nil, appError.NewAppError(400, "order value does not meet voucher minimum requirement")
		}
		// TODO: Validate usage limit, date, seller/category applicability

		discount := 0.0
		if voucher.DiscountType == "PERCENTAGE" {
			discount = totalAmount * float64(voucher.DiscountValue) / 100
			if discount > float64(voucher.MaxDiscountAmount) {
				discount = float64(voucher.MaxDiscountAmount)
			}
		} else if voucher.DiscountType == "FIXED" {
			discount = float64(voucher.DiscountValue)
		}

		if discount > float64(voucher.MaxDiscountAmount) {
			discount = float64(voucher.MaxDiscountAmount)
		}

		totalAmount -= discount
		if totalAmount < 0 {
			totalAmount = 0
		}

		orderVoucher = &model.OrderVoucher{
			Code:          voucher.Code,
			DiscountType:  voucher.DiscountType,
			DiscountValue: voucher.DiscountValue,
		}
	}

	if err != nil {
		// Rollback: release reserved stock
		_ = s.productClient.ReleaseStock(tempOrderID)
		return nil, fmt.Errorf("failed to create calculate fee request: %w", err)
	}
	if err != nil {
		// Rollback: release reserved stock
		_ = s.productClient.ReleaseStock(tempOrderID)
		return nil, fmt.Errorf("failed to get delivery fee: %w", err)
	}

	// 6. Create Order record
	orderStatus := "TO_PAY"
	if request.PaymentMethod == "COD" {
		orderStatus = "TO_CONFIRM"
	}

	order := &model.Order{
		Status:        orderStatus,
		PaymentMethod: request.PaymentMethod,
		PaymentStatus: "PENDING",
		User: model.User{
			ID:   user.ID,
			Name: user.Name,
		},
		Seller: model.User{
			ID:   seller.ID,
			Name: seller.Name,
		},
		Items:             orderItems,
		Voucher:           orderVoucher,
		Total:             totalAmount,
		DeliveryServiceID: request.DeliveryServiceID,
		ShippingAddress: model.OrderAddress{
			FullName:    request.ShippingAddress.FullName,
			Phone:       request.ShippingAddress.Phone,
			AddressLine: request.ShippingAddress.AddressLine,
			Ward:        request.ShippingAddress.Ward,
			District:    request.ShippingAddress.District,
			Province:    request.ShippingAddress.Province,
			Country:     request.ShippingAddress.Country,
			Latitude:    request.ShippingAddress.Latitude,
			Longitude:   request.ShippingAddress.Longitude,
			DistrictID:  request.ShippingAddress.DistrictID,
			WardCode:    request.ShippingAddress.WardCode,
		},
	}

	//handle delivery fee
	// calculate delivery fee
	calculateFeeRequest, err := s.GHNClient.CreateCalculateFeeRequest(*order, *sellerFetch)
	deliveryFeeResponse, err := s.GHNClient.CalculateFee(*calculateFeeRequest)
	if err != nil {
		// Rollback: release reserved stock
		_ = s.productClient.ReleaseStock(tempOrderID)
		return nil, fmt.Errorf("failed to get delivery fee: %w", err)
	}
	order.DeliveryFee = deliveryFeeResponse.Total

	if err := s.repo.CreateOrder(order); err != nil {
		// Rollback: release reserved stock
		_ = s.productClient.ReleaseStock(tempOrderID)
		return nil, err
	}

	// Update stock reservation with actual order ID
	// We need to release the temp reservation and create a new one with the actual order ID
	_ = s.productClient.ReleaseStock(tempOrderID)
	if err := s.productClient.ReserveStock(order.ID, reserveItems); err != nil {
		// Log error but don't fail the order creation since it's already created
		// In production, you might want to handle this differently
		fmt.Printf("Warning: failed to update stock reservation with order ID: %v\n", err)
	}

	// Delete purchased cart items
	for _, item := range cartItems {
		_ = s.cartRepo.DeleteCartItem(item.ID)
	}

	return &dto.CheckoutResponse{
		OrderID:     order.ID,
		TotalAmount: order.Total + float64(order.DeliveryFee),
		Status:      order.Status,
	}, nil
}

func (s *orderService) UpdateOrderPaymentStatus(ctx context.Context, orderID string, status string) error {
	order, err := s.repo.FindOrderByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return fmt.Errorf("order not found: %s", orderID)
	}

	order.PaymentStatus = status
	if status == "PAID" {
		order.Status = "PAID" // Or whatever the next status is, e.g. TO_SHIP
	} else if status == "FAILED" {
		order.Status = "CANCELLED"
	}

	return s.repo.UpdateOrder(order)
}

func (s *orderService) CreateCheckoutSession(order *model.Order, successURL, cancelURL string) (string, error) {
	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, item := range order.Items {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("vnd"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(item.ProductName + " - " + item.VariantName),
				},
				UnitAmount: stripe.Int64(int64(item.Price)),
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	// Add delivery fee as a separate line item
	if order.DeliveryFee > 0 {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("vnd"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String("Delivery Fee"),
				},
				UnitAmount: stripe.Int64(int64(order.DeliveryFee)),
			},
			Quantity: stripe.Int64(1),
		})
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems:  lineItems,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Metadata: map[string]string{
			"order_id": order.ID,
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	return sess.URL, nil
}

func (s *orderService) HandleCheckoutSessionCompleted(ctx context.Context, orderID, sessionID, paymentIntentID string) error {
	order, err := s.repo.FindOrderByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return fmt.Errorf("order not found: %s", orderID)
	}

	// Update payment status to PAID
	order.PaymentStatus = "PAID"
	// Update order status to TO_CONFIRM
	order.Status = "TO_CONFIRM"

	return s.repo.UpdateOrder(order)
}

func (s *orderService) HandlePaymentFailed(ctx context.Context, paymentIntentID string) error {
	// Find order by payment intent ID would require storing it
	// For now, we'll need to add a field to store payment intent ID
	// This is a simplified version that just logs the failure
	// In production, you'd want to find the order by payment intent ID
	// and update its status accordingly
	return nil
}

func (s *orderService) CreatePaymentForOrder(ctx context.Context, orderID string) (*dto.CreatePaymentResponse, error) {
	// Find the order
	order, err := s.repo.FindOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}

	// Create payment via Stripe
	paymentUrl, err := s.paymentClient.CreatePayment(order, s.clientURL+"/checkout/success", s.clientURL+"/checkout/failure")
	if err != nil {
		return nil, err
	}

	return &dto.CreatePaymentResponse{
		PaymentUrl: paymentUrl,
	}, nil
}

func (s *orderService) GetUserOrders(ctx context.Context, userID string, request dto.GetOrdersRequest) (*dto.GetOrdersResponse, error) {
	// Set defaults for pagination
	page := request.Page
	if page < 1 {
		page = 1
	}

	limit := request.Limit
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100 // max limit
	}

	// Query repository
	orders, totalCount, err := s.repo.FindOrdersByUser(userID, request.Status, page, limit, request.SortBy, request.SortOrder)
	if err != nil {
		return nil, err
	}

	// Map to DTOs
	orderDtos := make([]dto.OrderDto, len(orders))
	for i, order := range orders {
		orderDtos[i] = convertOrderToDto(order)
	}

	// Calculate total pages
	totalPages := int(totalCount) / limit
	if int(totalCount)%limit > 0 {
		totalPages++
	}

	return &dto.GetOrdersResponse{
		Orders:     orderDtos,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *orderService) GetSellerOrders(ctx context.Context, sellerID string, request dto.GetOrdersBySellerRequest) (*dto.GetOrdersResponse, error) {
	// Set defaults for pagination
	page := request.Page
	if page < 1 {
		page = 1
	}

	limit := request.Limit
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100 // max limit
	}

	// Query repository with filters and search
	orders, totalCount, err := s.repo.FindOrdersBySeller(
		sellerID,
		request.Status,
		request.PaymentMethod,
		request.PaymentStatus,
		request.Search,
		page,
		limit,
		request.SortBy,
		request.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	// Map to DTOs
	orderDtos := make([]dto.OrderDto, len(orders))
	for i, order := range orders {
		orderDtos[i] = convertOrderToDto(order)
	}

	// Calculate total pages
	totalPages := int(totalCount) / limit
	if int(totalCount)%limit > 0 {
		totalPages++
	}

	return &dto.GetOrdersResponse{
		Orders:     orderDtos,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, userID string, orderID string, request dto.UpdateOrderStatusRequest) error {

	//get old order
	oldOrder, err := s.repo.FindOrderByID(orderID)
	if err != nil {
		return err
	}
	if oldOrder == nil {
		return appError.NewAppError(404, "Order not found")
	}
	if oldOrder.Seller.ID != userID {
		return appError.NewAppError(403, "You are not authorized to update this order")
	}

	//fetch seller
	seller, err := s.userClient.GetUserByID(oldOrder.Seller.ID)
	if err != nil {
		return err
	}

	//validate current status
	switch oldOrder.Status {
	case "TO_CONFIRM":
		if request.Status != "TO_PICKUP" && request.Status != "CANCELLED" {
			return appError.NewAppError(400, "Invalid status")
		}
		if request.Status == "TO_PICKUP" {
			request, err := s.GHNClient.CreateRequest(*oldOrder, *seller)
			if err != nil {
				return err
			}
			deliveryCode, err := s.GHNClient.CreateOrder(request)
			if err != nil {
				return err
			}
			oldOrder.DeliveryCode = deliveryCode
		}
	case "TO_PICKUP":
		if request.Status != "SHIPPING" && request.Status != "CANCELLED" {
			return appError.NewAppError(400, "Invalid status")
		}
	case "SHIPPING":
		if request.Status != "COMPLETED" && request.Status != "RETURNED" {
			return appError.NewAppError(400, "Invalid status")
		}
	default:
		return appError.NewAppError(409, "Current status of this order can't be updated"+oldOrder.Status)
	}

	//update in database
	oldOrder.Status = request.Status
	return s.repo.UpdateOrder(oldOrder)
}

// Helper function to convert Order model to OrderDto
func convertOrderToDto(order *model.Order) dto.OrderDto {
	// Convert items
	itemDtos := make([]dto.OrderItemDto, len(order.Items))
	for i, item := range order.Items {
		itemDtos[i] = dto.OrderItemDto{
			ProductID:   item.ProductID,
			VariantID:   item.VariantID,
			ProductName: item.ProductName,
			VariantName: item.VariantName,
			SKU:         item.SKU,
			Price:       item.Price,
			Image:       item.Image,
			Quantity:    item.Quantity,
		}
	}

	// Convert voucher if present
	var voucherDto *dto.OrderVoucherDto
	if order.Voucher != nil {
		voucherDto = &dto.OrderVoucherDto{
			Code:                   order.Voucher.Code,
			DiscountType:           order.Voucher.DiscountType,
			DiscountValue:          order.Voucher.DiscountValue,
			MaxDiscountValue:       order.Voucher.MaxDiscountValue,
			MinOrderValue:          order.Voucher.MinOrderValue,
			ApplyScope:             order.Voucher.ApplyScope,
			ApplySellerCategoryIds: order.Voucher.ApplySellerCategoryIds,
		}
	}

	return dto.OrderDto{
		ID:     order.ID,
		Status: order.Status,
		User: dto.UserDto{
			ID:       order.User.ID,
			Username: order.User.Username,
			Name:     order.User.Name,
			Email:    order.User.Email,
		},
		PaymentMethod: order.PaymentMethod,
		PaymentStatus: order.PaymentStatus,
		Seller: dto.UserDto{
			ID:       order.Seller.ID,
			Username: order.Seller.Username,
			Name:     order.Seller.Name,
			Email:    order.Seller.Email,
		},
		Items:     itemDtos,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		Voucher:   voucherDto,
		Total:     order.Total,
		Phone:     order.Phone,
		ShippingAddress: dto.OrderAddressDto{
			FullName:    order.ShippingAddress.FullName,
			Phone:       order.ShippingAddress.Phone,
			AddressLine: order.ShippingAddress.AddressLine,
			Ward:        order.ShippingAddress.Ward,
			District:    order.ShippingAddress.District,
			Province:    order.ShippingAddress.Province,
			Country:     order.ShippingAddress.Country,
			Latitude:    order.ShippingAddress.Latitude,
			Longitude:   order.ShippingAddress.Longitude,
		},
		DeliveryCode: order.DeliveryCode,
		ItemCount:    len(order.Items),
	}
}
