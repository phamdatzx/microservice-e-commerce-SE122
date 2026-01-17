package dto

type UpdateProductCountRequest struct {
	SellerID  string `json:"seller_id" binding:"required"`
	Operation string `json:"operation" binding:"required,oneof=increment decrement"`
}

type UpdateProductCountResponse struct {
	ProductCount int `json:"product_count"`
}
