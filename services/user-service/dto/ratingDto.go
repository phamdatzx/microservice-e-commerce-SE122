package dto

type UpdateRatingRequest struct {
	SellerID  string  `json:"seller_id" binding:"required"`
	Star      float64 `json:"star" binding:"required,min=1,max=5"`
	Operation string  `json:"operation" binding:"required,oneof=create update delete"`
	OldStar   float64 `json:"old_star"` // Only required for update operation
}

type UpdateRatingResponse struct {
	RatingCount   int     `json:"rating_count"`
	RatingAverage float64 `json:"rating_average"`
}
