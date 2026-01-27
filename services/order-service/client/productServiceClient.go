package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"order-service/dto"
	"os"
	"time"
)

// ProductServiceClient handles communication with product-service
type ProductServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewProductServiceClient creates a new product service client
func NewProductServiceClient() *ProductServiceClient {
	baseURL := os.Getenv("PRODUCT_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081" // Default for local development
	}

	return &ProductServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetVariantsByIdsRequest represents the request to product-service
type GetVariantsByIdsRequest struct {
	VariantIDs []string `json:"variant_ids"`
}

// GetVariantsByIdsResponse represents the response from product-service
type GetVariantsByIdsResponse struct {
	Variants []dto.ProductVariantDto `json:"variants"`
}

// GetVariantsByIds calls product-service to get variant details
func (c *ProductServiceClient) GetVariantsByIds(variantIDs []string) ([]dto.ProductVariantDto, error) {
	if len(variantIDs) == 0 {
		return []dto.ProductVariantDto{}, nil
	}

	// Prepare request
	requestBody := GetVariantsByIdsRequest{
		VariantIDs: variantIDs,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request
	url := fmt.Sprintf("%s/api/product/public/variants/batch", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call product-service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("product-service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response GetVariantsByIdsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Variants, nil
}

func (c *ProductServiceClient) GetVoucherByID(voucherID string) (*dto.VoucherResponse, error) {
	url := fmt.Sprintf("%s/api/product/vouchers/%s", c.baseURL, voucherID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call product-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product-service returned status %d", resp.StatusCode)
	}

	var response dto.VoucherResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// ReserveStockItem represents a single item to reserve
type ReserveStockItem struct {
	VariantID string `json:"variant_id"`
	Quantity  int    `json:"quantity"`
}

// ReserveStockRequest represents the request to reserve stock
type ReserveStockRequest struct {
	OrderID string             `json:"order_id"`
	Items   []ReserveStockItem `json:"items"`
}

// ReleaseStockRequest represents the request to release stock
type ReleaseStockRequest struct {
	OrderID string `json:"order_id"`
}

// ReserveStock calls product-service to reserve stock for an order
func (c *ProductServiceClient) ReserveStock(orderID string, items []ReserveStockItem) error {
	// Prepare request
	requestBody := ReserveStockRequest{
		OrderID: orderID,
		Items:   items,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request
	url := fmt.Sprintf("%s/api/product/stock/reserve", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call product-service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("product-service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ReleaseStock calls product-service to release stock for an order
func (c *ProductServiceClient) ReleaseStock(orderID string) error {
	// Prepare request
	requestBody := ReleaseStockRequest{
		OrderID: orderID,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request
	url := fmt.Sprintf("%s/api/product/stock/release", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call product-service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("product-service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// UseVoucher calls product-service to use a voucher
func (c *ProductServiceClient) UseVoucher(userID, voucherID string) (*dto.UseVoucherResponse, error) {
	// Prepare request
	requestBody := dto.UseVoucherRequest{
		UserID:    userID,
		VoucherID: voucherID,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request
	url := fmt.Sprintf("%s/api/product/vouchers/use", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call product-service: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var response dto.UseVoucherResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return &response, fmt.Errorf("product-service returned status %d: %s", resp.StatusCode, response.Message)
	}

	return &response, nil
}

// ProductStatusInfo represents status information for a single product
type ProductStatusInfo struct {
	ProductID  string `json:"product_id"`
	IsReported bool   `json:"is_reported"`
	IsRated    bool   `json:"is_rated"`
}

// CheckProductStatusRequest represents the request to check product status
type CheckProductStatusRequest struct {
	UserID     string   `json:"user_id"`
	ProductIDs []string `json:"product_ids"`
}

// CheckProductStatusResponse represents the response for checking product status
type CheckProductStatusResponse struct {
	Products []ProductStatusInfo `json:"products"`
}

// CheckProductsStatus calls product-service to check if products are rated/reported by user
func (c *ProductServiceClient) CheckProductsStatus(userID string, productIDs []string) (map[string]ProductStatusInfo, error) {
	if len(productIDs) == 0 {
		return make(map[string]ProductStatusInfo), nil
	}

	// Prepare request
	requestBody := CheckProductStatusRequest{
		UserID:     userID,
		ProductIDs: productIDs,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request
	url := fmt.Sprintf("%s/api/product/public/check-status", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call product-service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("product-service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response CheckProductStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to map for easy lookup
	statusMap := make(map[string]ProductStatusInfo)
	for _, product := range response.Products {
		statusMap[product.ProductID] = product
	}

	return statusMap, nil
}
