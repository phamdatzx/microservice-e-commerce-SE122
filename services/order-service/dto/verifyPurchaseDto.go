package dto

type VerifyPurchaseRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	ProductID string `json:"product_id" binding:"required"`
	VariantID string `json:"variant_id" binding:"required"`
}

type VerifyPurchaseResponse struct {
	HasPurchased bool `json:"has_purchased"`
}
