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

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
)

type OrderService interface {
	Checkout(userID string, request dto.CheckoutRequest) (*dto.CheckoutResponse, error)
	UpdateOrderPaymentStatus(ctx context.Context, orderID string, status string) error
	CreateCheckoutSession(order *model.Order, successURL, cancelURL string) (string, error)
}

type orderService struct {
	repo          repository.OrderRepository
	cartRepo      repository.CartRepository
	productClient *client.ProductServiceClient
	userClient    *client.UserServiceClient
	paymentClient payment.PaymentClient
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	productClient *client.ProductServiceClient,
	userClient *client.UserServiceClient,
	paymentClient payment.PaymentClient,
) OrderService {
	return &orderService{
		repo:          orderRepo,
		cartRepo:      cartRepo,
		productClient: productClient,
		userClient:    userClient,
		paymentClient: paymentClient,
	}
}

func (s *orderService) Checkout(userID string, request dto.CheckoutRequest) (*dto.CheckoutResponse, error) {
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

	// 4. Fetch Variant details and Verify stock
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

	// 5. Apply Voucher
	var orderVoucher *model.OrderVoucher
	if request.VoucherID != "" {
		voucher, err := s.productClient.GetVoucherByID(request.VoucherID)
		if err != nil {
			return nil, fmt.Errorf("failed to get voucher: %w", err)
		}

		// Validate voucher (basic validation, more complex logic might be needed)
		if voucher.Status != "active" {
			return nil, appError.NewAppError(400, "voucher is not active")
		}
		if totalAmount < float64(voucher.MinOrderValue) {
			return nil, appError.NewAppError(400, "order value does not meet voucher minimum requirement")
		}
		// TODO: Validate usage limit, date, seller/category applicability

		discount := 0.0
		if voucher.DiscountType == "percentage" {
			discount = totalAmount * float64(voucher.DiscountValue) / 100
			if discount > float64(voucher.MaxDiscountAmount) {
				discount = float64(voucher.MaxDiscountAmount)
			}
		} else if voucher.DiscountType == "fixed_amount" {
			discount = float64(voucher.DiscountValue)
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

	// 6. Create Order record
	orderStatus := "TO_PAY"
	if request.PaymentMethod == "COD" {
		orderStatus = "TO_CONFIRM"
	}

	
	order := &model.Order{
		Status: orderStatus,
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
		Items:   orderItems,
		Voucher: orderVoucher,
		Total:   totalAmount,
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
		},
	}

	if err := s.repo.CreateOrder(order); err != nil {
		return nil, err
	}

	// 7. Handle Payment if needed
	var clientSecret string
	var paymentUrl string
	if request.PaymentMethod == "STRIPE" {
		var err error
		paymentUrl, err = s.CreateCheckoutSession(order, "http://localhost:3000/checkout/success", "http://localhost:3000/checkout/failure")
		if err != nil {
			return nil, err
		}
	}

	// 8. Delete purchased cart items
	for _, item := range cartItems {
		_ = s.cartRepo.DeleteCartItem(item.ID)
	}

	return &dto.CheckoutResponse{
		OrderID:      order.ID,
		TotalAmount:  order.Total,
		Status:       order.Status,
		ClientSecret: clientSecret,
		PaymentUrl:   paymentUrl,
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
				Currency: stripe.String("usd"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(item.ProductName + " - " + item.VariantName),
				},
				UnitAmount: stripe.Int64(int64(item.Price * 100)),
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
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
