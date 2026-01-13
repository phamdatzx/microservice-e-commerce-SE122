package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type OrderServiceClient struct {
	baseURL string
	client  *http.Client
}

func NewOrderServiceClient() *OrderServiceClient {
	return &OrderServiceClient{
		baseURL: os.Getenv("ORDER_SERVICE_URL"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type VerifyPurchaseRequest struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id"`
}

type VerifyPurchaseResponse struct {
	HasPurchased bool `json:"has_purchased"`
}

func (c *OrderServiceClient) VerifyVariantPurchase(userID, productID, variantID string) (bool, error) {
	url := fmt.Sprintf("%s/api/order/verify-purchase", c.baseURL)

	requestBody := VerifyPurchaseRequest{
		UserID:    userID,
		ProductID: productID,
		VariantID: variantID,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return false, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("verify purchase failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response VerifyPurchaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.HasPurchased, nil
}
