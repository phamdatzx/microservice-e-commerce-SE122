package dto

// ChatBotRequest represents the request body for chatbot API
type ChatBotRequest struct {
	Type      string `json:"type" binding:"required,oneof=general product seller"` // general, product, seller
	SellerID  string `json:"seller_id"`
	ProductID string `json:"product_id"`
	Question  string `json:"question" binding:"required"`
}

// ChatBotResponse represents the response from chatbot API
type ChatBotResponse struct {
	Answer string `json:"answer"`
}
