package dto

import "product-service/model"

type CartVariantDto struct {
	ProductName       string        `json:"product_name"`
	SellerID          string        `json:"seller_id"`
	SellerCategoryIds []string      `json:"seller_category_ids"`
	Variant           model.Variant `json:"variant"`
}

// GetVariantsByIdsRequest represents the request body for getting variants by IDs
type GetVariantsByIdsRequest struct {
	VariantIDs []string `json:"variant_ids" binding:"required,min=1"`
}

// GetVariantsByIdsResponse represents the response for getting variants by IDs
type GetVariantsByIdsResponse struct {
	Variants []CartVariantDto `json:"variants"`
}
