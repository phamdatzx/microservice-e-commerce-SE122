package service

import (
	"fmt"
	"product-service/client"
	"product-service/dto"
	appError "product-service/error"
	"product-service/model"
	"product-service/repository"
	"strings"
)

type ChatBotService interface {
	AskQuestion(request dto.ChatBotRequest) (string, error)
}

type chatBotService struct {
	chatBotClient      *client.ChatBotClient
	productRepo        repository.ProductRepository
	categoryRepo       repository.CategoryRepository
	sellerCategoryRepo repository.SellerCategoryRepository
	voucherRepo        repository.VoucherRepository
	userClient         *client.UserServiceClient
}

func NewChatBotService(
	chatBotClient *client.ChatBotClient,
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	sellerCategoryRepo repository.SellerCategoryRepository,
	voucherRepo repository.VoucherRepository,
	userClient *client.UserServiceClient,
) ChatBotService {
	return &chatBotService{
		chatBotClient:      chatBotClient,
		productRepo:        productRepo,
		categoryRepo:       categoryRepo,
		sellerCategoryRepo: sellerCategoryRepo,
		voucherRepo:        voucherRepo,
		userClient:         userClient,
	}
}

func (s *chatBotService) AskQuestion(request dto.ChatBotRequest) (string, error) {
	var context string
	var err error

	switch request.Type {
	case "general":
		context, err = s.buildGeneralContext()
	case "product":
		if request.ProductID == "" {
			return "", appError.NewAppError(400, "product_id is required for product type questions")
		}
		context, err = s.buildProductContext(request.ProductID)
	case "seller":
		if request.SellerID == "" {
			return "", appError.NewAppError(400, "seller_id is required for seller type questions")
		}
		context, err = s.buildSellerContext(request.SellerID)
	default:
		return "", appError.NewAppError(400, "invalid question type")
	}

	if err != nil {
		return "", err
	}

	// Send question to AI model
	answer, err := s.chatBotClient.SendMessage(request.Question, context)
	if err != nil {
		return "", appError.NewAppErrorWithErr(500, "failed to get response from AI model", err)
	}

	return answer, nil
}

// buildGeneralContext builds context for general questions
// Includes: categories, popular products, popular sellers
func (s *chatBotService) buildGeneralContext() (string, error) {
	var contextBuilder strings.Builder

	// Get all categories
	categories, err := s.categoryRepo.FindAll()
	if err == nil && len(categories) > 0 {
		contextBuilder.WriteString("=== DANH MỤC SẢN PHẨM ===\n")
		for _, cat := range categories {
			contextBuilder.WriteString(fmt.Sprintf("- %s (ID: %s)\n", cat.Name, cat.ID))
		}
		contextBuilder.WriteString("\n")
	}

	// Get popular products (top sold products)
	// Using search with sort by sold_count
	products, _, err := s.productRepo.SearchProducts(
		map[string]interface{}{"is_disabled": map[string]interface{}{"$ne": true}},
		0, 10, false, "sold_count", -1,
	)
	if err == nil && len(products) > 0 {
		contextBuilder.WriteString("=== SẢN PHẨM PHỔ BIẾN ===\n")
		for _, p := range products {
			contextBuilder.WriteString(s.formatProductInfo(&p))
		}
		contextBuilder.WriteString("\n")
	}

	return contextBuilder.String(), nil
}

// buildProductContext builds context for product-specific questions
// Includes: product details, seller info, seller categories, vouchers
func (s *chatBotService) buildProductContext(productID string) (string, error) {
	var contextBuilder strings.Builder

	// Get product info
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return "", appError.NewAppErrorWithErr(404, "product not found", err)
	}

	contextBuilder.WriteString("=== THÔNG TIN SẢN PHẨM ===\n")
	contextBuilder.WriteString(s.formatProductDetailInfo(product))
	contextBuilder.WriteString("\n")

	// Get seller info
	if product.SellerID != "" {
		seller, err := s.userClient.GetUserByID(product.SellerID)
		if err == nil && seller != nil {
			contextBuilder.WriteString("=== THÔNG TIN NGƯỜI BÁN ===\n")
			contextBuilder.WriteString(fmt.Sprintf("- Tên: %s\n", seller.Name))
			contextBuilder.WriteString("\n")
		}

		// Get seller categories
		sellerCategories, err := s.sellerCategoryRepo.FindBySellerID(product.SellerID)
		if err == nil && len(sellerCategories) > 0 {
			contextBuilder.WriteString("=== DANH MỤC CỦA NGƯỜI BÁN ===\n")
			for _, cat := range sellerCategories {
				contextBuilder.WriteString(fmt.Sprintf("- %s\n", cat.Name))
			}
			contextBuilder.WriteString("\n")
		}

		// Get seller vouchers
		vouchers, err := s.voucherRepo.FindBySellerID(product.SellerID)
		if err == nil && len(vouchers) > 0 {
			contextBuilder.WriteString("=== VOUCHER CỦA CỬA HÀNG ===\n")
			for _, v := range vouchers {
				if v.Status == "ACTIVE" {
					contextBuilder.WriteString(s.formatVoucherInfo(&v))
				}
			}
			contextBuilder.WriteString("\n")
		}
	}

	return contextBuilder.String(), nil
}

// buildSellerContext builds context for seller-specific questions
// Includes: seller categories, vouchers
func (s *chatBotService) buildSellerContext(sellerID string) (string, error) {
	var contextBuilder strings.Builder

	// Get seller info
	seller, err := s.userClient.GetUserByID(sellerID)
	if err == nil && seller != nil {
		contextBuilder.WriteString("=== THÔNG TIN CỬA HÀNG ===\n")
		contextBuilder.WriteString(fmt.Sprintf("- Tên: %s\n", seller.Name))
		contextBuilder.WriteString("\n")
	}

	// Get seller categories
	sellerCategories, err := s.sellerCategoryRepo.FindBySellerID(sellerID)
	if err == nil && len(sellerCategories) > 0 {
		contextBuilder.WriteString("=== DANH MỤC SẢN PHẨM CỦA CỬA HÀNG ===\n")
		for _, cat := range sellerCategories {
			contextBuilder.WriteString(fmt.Sprintf("- %s (%d sản phẩm)\n", cat.Name, cat.ProductCount))
		}
		contextBuilder.WriteString("\n")
	}

	// Get seller vouchers
	vouchers, err := s.voucherRepo.FindBySellerID(sellerID)
	if err == nil && len(vouchers) > 0 {
		contextBuilder.WriteString("=== VOUCHER CỦA CỬA HÀNG ===\n")
		for _, v := range vouchers {
			if v.Status == "ACTIVE" {
				contextBuilder.WriteString(s.formatVoucherInfo(&v))
			}
		}
		contextBuilder.WriteString("\n")
	}

	return contextBuilder.String(), nil
}

func (s *chatBotService) formatProductInfo(p *model.Product) string {
	return fmt.Sprintf("- %s | Giá: %d-%d VND | Đánh giá: %.1f/5 | Đã bán: %d\n",
		p.Name, p.Price.Min, p.Price.Max, p.Rating, p.SoldCount)
}

func (s *chatBotService) formatProductDetailInfo(p *model.Product) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Tên: %s\n", p.Name))
	sb.WriteString(fmt.Sprintf("Mô tả: %s\n", p.Description))
	sb.WriteString(fmt.Sprintf("Giá: %d - %d VND\n", p.Price.Min, p.Price.Max))
	sb.WriteString(fmt.Sprintf("Đánh giá: %.1f/5 (%d đánh giá)\n", p.Rating, p.RateCount))
	sb.WriteString(fmt.Sprintf("Đã bán: %d\n", p.SoldCount))
	sb.WriteString(fmt.Sprintf("Tình trạng: %s\n", p.Status))
	sb.WriteString(fmt.Sprintf("Tồn kho: %d\n", p.Stock))

	// Add variant info
	if len(p.Variants) > 0 {
		sb.WriteString("Các biến thể:\n")
		for _, v := range p.Variants {
			optionStr := ""
			for key, val := range v.Options {
				optionStr += fmt.Sprintf("%s: %s, ", key, val)
			}
			sb.WriteString(fmt.Sprintf("  - %s | Giá: %d VND | Còn: %d\n", strings.TrimSuffix(optionStr, ", "), v.Price, v.Stock))
		}
	}

	return sb.String()
}

func (s *chatBotService) formatVoucherInfo(v *model.Voucher) string {
	discountStr := ""
	if v.DiscountType == "PERCENTAGE" {
		discountStr = fmt.Sprintf("Giảm %d%%", v.DiscountValue)
		if v.MaxDiscountValue != nil {
			discountStr += fmt.Sprintf(" (tối đa %d VND)", *v.MaxDiscountValue)
		}
	} else {
		discountStr = fmt.Sprintf("Giảm %d VND", v.DiscountValue)
	}

	return fmt.Sprintf("- %s: %s | Đơn tối thiểu: %d VND | Code: %s\n",
		v.Name, discountStr, v.MinOrderValue, v.Code)
}
