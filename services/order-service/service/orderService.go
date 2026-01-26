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
	"slices"
	"time"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
)

type OrderService interface {
	Checkout(userID string, request dto.CheckoutRequest) (*dto.CheckoutResponse, error)
	UpdateOrderPaymentStatus(ctx context.Context, orderID string, status string) error
	CreateCheckoutSession(ctx context.Context, orderID string) (*dto.CreatePaymentResponse, error)
	HandleCheckoutSessionCompleted(ctx context.Context, orderID, sessionID, paymentIntentID string) error
	HandlePaymentFailed(ctx context.Context, paymentIntentID string) error
	GetUserOrders(ctx context.Context, userID string, request dto.GetOrdersRequest) (*dto.GetOrdersResponse, error)
	GetSellerOrders(ctx context.Context, sellerID string, request dto.GetOrdersBySellerRequest) (*dto.GetOrdersResponse, error)
	UpdateOrderStatus(ctx context.Context, userID string, orderID string, request dto.UpdateOrderStatusRequest) error
	VerifyVariantPurchase(userID, productID, variantID string) (bool, error)
	GetSellerStatistics(ctx context.Context, sellerID string, request dto.GetSellerStatisticsRequest) (*dto.GetSellerStatisticsResponse, error)
	InstantCheckout(userID string, request dto.InstantCheckoutRequest) (*dto.CheckoutResponse, error)
	ApplyVoucher(voucherId string, totalAmount float64, sellerId string, variants []dto.ProductVariantDto, userId string) (float64, *model.OrderVoucher, error)
}

type orderService struct {
	repo               repository.OrderRepository
	cartRepo           repository.CartRepository
	productClient      *client.ProductServiceClient
	userClient         *client.UserServiceClient
	paymentClient      payment.PaymentClient
	GHNClient          *client.GHNClient
	notificationClient *client.NotificationServiceClient
	clientURL          string
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	productClient *client.ProductServiceClient,
	userClient *client.UserServiceClient,
	paymentClient payment.PaymentClient,
	GHNClient *client.GHNClient,
	notificationClient *client.NotificationServiceClient,
) OrderService {
	return &orderService{
		repo:               orderRepo,
		cartRepo:           cartRepo,
		productClient:      productClient,
		userClient:         userClient,
		paymentClient:      paymentClient,
		GHNClient:          GHNClient,
		notificationClient: notificationClient,
		clientURL:          os.Getenv("CLIENT_URL"),
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
	var saveOrderVoucher *model.OrderVoucher
	if request.VoucherID != "" {
		newTotalAmount, orderVoucher, err := s.ApplyVoucher(request.VoucherID, totalAmount, sellerID, variants, userID)
		if err != nil {
			// Rollback: release reserved stock
			_ = s.productClient.ReleaseStock(tempOrderID)
			return nil, err
		}

		totalAmount = newTotalAmount
		saveOrderVoucher = orderVoucher
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
		Voucher:           saveOrderVoucher,
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

	// Send notification to seller for COD orders
	if request.PaymentMethod == "COD" {
		orderData := map[string]interface{}{
			"orderId":      order.ID,
			"total":        order.Total + float64(order.DeliveryFee),
			"itemCount":    len(order.Items),
			"customerName": user.Name,
		}

		err = s.notificationClient.CreateNotification(client.CreateNotificationRequest{
			UserID:  sellerID,
			Type:    "order",
			Title:   "New COD Order",
			Message: fmt.Sprintf("You have a new COD order from %s", user.Name),
			Data:    orderData,
		})

		if err != nil {
			// Log error but don't fail the checkout
			fmt.Printf("Warning: failed to send notification to seller: %v\n", err)
		}
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

func (s *orderService) CreateCheckoutSession(ctx context.Context, orderID string) (*dto.CreatePaymentResponse, error) {

	// Find the order
	order, err := s.repo.FindOrderByID(orderID)
	if err != nil {
		return nil,  err
	}
	if order == nil {
		return nil,  appError.NewAppError(404 ,"order not found")
	}

	var lineItems []*stripe.CheckoutSessionLineItemParams

	// Add order total (excluding delivery fee) as a line item
	if order.Total > 0 {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("vnd"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String("Order Total"),
				},
				UnitAmount: stripe.Int64(int64(order.Total)),
			},
			Quantity: stripe.Int64(1),
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
		SuccessURL: stripe.String(s.clientURL+"/checkout/success"),
		CancelURL:  stripe.String(s.clientURL+"/checkout/failure"),
		Metadata: map[string]string{
			"order_id": order.ID,
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, appError.NewAppError(500, "failed to create checkout session")
	}

	return &dto.CreatePaymentResponse{
		PaymentUrl: sess.URL,
	}, nil
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

	err = s.repo.UpdateOrder(order)
	if err != nil {
		return err
	}

	// Send notification to seller about the paid order
	orderData := map[string]interface{}{
		"orderId":      order.ID,
		"total":        order.Total + float64(order.DeliveryFee),
		"itemCount":    len(order.Items),
		"customerName": order.User.Name,
		"paymentMethod": order.PaymentMethod,
	}

	err = s.notificationClient.CreateNotification(client.CreateNotificationRequest{
		UserID:  order.Seller.ID,
		Type:    "order",
		Title:   "New Paid Order",
		Message: fmt.Sprintf("You have a new paid order from %s", order.User.Name),
		Data:    orderData,
	})

	if err != nil {
		// Log error but don't fail the order update
		fmt.Printf("Warning: failed to send notification to seller: %v\n", err)
	}

	return nil
}

func (s *orderService) HandlePaymentFailed(ctx context.Context, paymentIntentID string) error {
	// Find order by payment intent ID would require storing it
	// For now, we'll need to add a field to store payment intent ID
	// This is a simplified version that just logs the failure
	// In production, you'd want to find the order by payment intent ID
	// and update its status accordingly
	return nil
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
	if oldOrder.Seller.ID != userID && oldOrder.User.ID != userID {
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
		if request.Status != "SHIPPING" && request.Status != "CANCELLED" && request.Status != "COMPLETED" {
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

func (s *orderService) VerifyVariantPurchase(userID, productID, variantID string) (bool, error) {
	return s.repo.VerifyVariantPurchase(userID, productID, variantID)
}

func (s *orderService) GetSellerStatistics(ctx context.Context, sellerID string, request dto.GetSellerStatisticsRequest) (*dto.GetSellerStatisticsResponse, error) {
	// Validate time range
	if request.From.After(request.To) {
		return nil, appError.NewAppError(400, "invalid time range: 'from' must be before 'to'")
	}

	// Validate type parameter
	groupBy := request.Type
	if groupBy != "" && groupBy != "day" && groupBy != "month" {
		return nil, appError.NewAppError(400, "invalid type parameter: must be 'day' or 'month'")
	}

	// Get statistics from repository
	orderCount, totalRevenue, breakdownData, err := s.repo.GetSellerStatistics(sellerID, request.From, request.To, groupBy)
	if err != nil {
		return nil, err
	}

	response := &dto.GetSellerStatisticsResponse{
		OrderCount:   orderCount,
		TotalRevenue: totalRevenue,
	}

	// If breakdown is requested, convert the data
	if groupBy != "" && breakdownData != nil {
		breakdown := make([]dto.PeriodStatistics, 0, len(breakdownData))
		for _, item := range breakdownData {
			period := item["_id"].(string)
			count := int(item["count"].(int32))
			revenue := item["revenue"].(float64)
			
			breakdown = append(breakdown, dto.PeriodStatistics{
				Period:     period,
				OrderCount: count,
				Revenue:    revenue,
			})
		}
		response.Breakdown = breakdown
	}

	return response, nil
}

func (s *orderService) InstantCheckout(userID string, request dto.InstantCheckoutRequest) (*dto.CheckoutResponse, error) {
	// Validate all items from same seller
	if len(request.Items) == 0 {
		return nil, appError.NewAppError(400, "no items provided")
	}

	sellerID := request.Items[0].SellerID
	for _, item := range request.Items {
		if item.SellerID != sellerID {
			return nil, appError.NewAppError(400, "all items must be from the same seller")
		}
	}

	// Fetch seller info for delivery
	sellerFetch, err := s.userClient.GetUserByID(sellerID)
	if err != nil {
		return nil, err
	}

	// Fetch User and Seller info
	user, err := s.userClient.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	seller, err := s.userClient.GetUserByID(sellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller info: %w", err)
	}

	// Fetch Variant details and validate stock
	var orderItems []model.OrderItem
	var totalAmount float64
	var variantIDs []string

	for _, item := range request.Items {
		variantIDs = append(variantIDs, item.VariantID)
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

	for _, item := range request.Items {
		variantInfo, exists := variantMap[item.VariantID]
		if !exists {
			return nil, appError.NewAppError(404, fmt.Sprintf("variant %s not found", item.VariantID))
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
			ProductID:   item.ProductID,
			ProductName: variantInfo.ProductName,
			VariantID:   item.VariantID,
			VariantName: variantName,
			SKU:         variantInfo.Variant.SKU,
			Price:       variantInfo.Variant.Price,
			Quantity:    item.Quantity,
			Image:       variantInfo.Variant.Image,
		})
	}

	// Reserve stock for all items
	var reserveItems []client.ReserveStockItem
	for _, item := range request.Items {
		reserveItems = append(reserveItems, client.ReserveStockItem{
			VariantID: item.VariantID,
			Quantity:  item.Quantity,
		})
	}

	// Create a temporary order ID for stock reservation
	tempOrderID := fmt.Sprintf("temp_%s_%d", userID, time.Now().UnixNano())

	if err := s.productClient.ReserveStock(tempOrderID, reserveItems); err != nil {
		return nil, fmt.Errorf("failed to reserve stock: %w", err)
	}

	// Apply Voucher
	var saveOrderVoucher *model.OrderVoucher
	if request.VoucherID != "" {
		newTotalAmount, orderVoucher, err := s.ApplyVoucher(request.VoucherID, totalAmount, sellerID,variants,userID)
		if err != nil {
			// Rollback: release reserved stock
			_ = s.productClient.ReleaseStock(tempOrderID)
			return nil, err
		}

		totalAmount = newTotalAmount
		saveOrderVoucher = orderVoucher
	}

	// Create Order record
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
		Voucher:           saveOrderVoucher,
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

	// Calculate delivery fee
	calculateFeeRequest, err := s.GHNClient.CreateCalculateFeeRequest(*order, *sellerFetch)
	if err != nil {
		// Rollback: release reserved stock
		_ = s.productClient.ReleaseStock(tempOrderID)
		return nil, fmt.Errorf("failed to create calculate fee request: %w", err)
	}
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
	_ = s.productClient.ReleaseStock(tempOrderID)
	if err := s.productClient.ReserveStock(order.ID, reserveItems); err != nil {
		// Log error but don't fail the order creation since it's already created
		fmt.Printf("Warning: failed to update stock reservation with order ID: %v\n", err)
	}

	// Send notification to seller for COD orders
	if request.PaymentMethod == "COD" {
		orderData := map[string]interface{}{
			"orderId":      order.ID,
			"total":        order.Total + float64(order.DeliveryFee),
			"itemCount":    len(order.Items),
			"customerName": user.Name,
		}

		err = s.notificationClient.CreateNotification(client.CreateNotificationRequest{
			UserID:  sellerID,
			Type:    "order",
			Title:   "New COD Order",
			Message: fmt.Sprintf("You have a new COD order from %s", user.Name),
			Data:    orderData,
		})

		if err != nil {
			// Log error but don't fail the checkout
			fmt.Printf("Warning: failed to send notification to seller: %v\n", err)
		}
	}

	return &dto.CheckoutResponse{
		OrderID:     order.ID,
		TotalAmount: order.Total + float64(order.DeliveryFee),
		Status:      order.Status,
	}, nil
}

func (s *orderService) ApplyVoucher(voucherId string, totalAmount float64, sellerId string, variants []dto.ProductVariantDto, userId string) (float64, *model.OrderVoucher, error) {

	voucher, err := s.productClient.GetVoucherByID(voucherId)
	if err != nil {
		return 0, nil, err
	}

	// Voucher should belong to seller
	if voucher.SellerID != sellerId {
		return 0, nil, appError.NewAppError(400, "voucher does not belong to seller")
	}

	// Validate status
	if voucher.Status != "ACTIVE" {
		return 0, nil, appError.NewAppError(400, "voucher is not active")	
	}

	// Validate min order value
	if totalAmount < float64(voucher.MinOrderValue) {
		return 0, nil, appError.NewAppError(400, "order value does not meet voucher minimum requirement")
	}

	//validate time 
	if time.Now().Before(voucher.StartDate) || time.Now().After(voucher.EndDate) {
		return 0, nil, appError.NewAppError(400, "voucher time has ended or not started")
	}
	//validate scope
	if voucher.ApplyScope == "CATEGORY" {
		if voucher.ApplySellerCategoryIds == nil {
			return 0, nil, appError.NewAppError(400, "voucher apply scope is category but no category id")
		}
		
		for _, variant := range variants {
			if !HasAny(voucher.ApplySellerCategoryIds, variant.CategoryIds) {
				return 0, nil, appError.NewAppError(400, "voucher apply scope is category but no category id")
			}
		}
	}

	//call use voucher api
	useVoucherResponse, err := s.productClient.UseVoucher(voucherId, userId)
	if err != nil {
		return 0, nil, appError.NewAppError(400, useVoucherResponse.Message)
	}

	// Calculate discount
	var discount float64
	if voucher.DiscountType == "PERCENTAGE" {
		discount = totalAmount * float64(voucher.DiscountValue) / 100
		if discount > float64(voucher.MaxDiscountValue) {
			discount = float64(voucher.MaxDiscountValue)
		}
	} else if voucher.DiscountType == "FIXED" {
		discount = float64(voucher.DiscountValue)
	}

	if voucher.MaxDiscountValue > 0 && discount > float64(voucher.MaxDiscountValue) {
		discount = float64(voucher.MaxDiscountValue)
	}

	totalAmount -= discount
	if totalAmount < 0 {
		totalAmount = 0
	}

	orderVoucher := &model.OrderVoucher{
		Code:          voucher.Code,
		DiscountType:  voucher.DiscountType,
		DiscountValue: voucher.DiscountValue,
	}

	return totalAmount, orderVoucher, nil
}
	
func HasAny(a, b []string) bool {
	for _, x := range a {
		if slices.Contains(b, x) {
			return true
		}
	}
	return false
}
