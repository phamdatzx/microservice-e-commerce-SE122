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

type UserServiceClient struct {
	baseURL string
	client  *http.Client
}

func NewUserServiceClient() *UserServiceClient {
	return &UserServiceClient{
		baseURL: os.Getenv("USER_SERVICE_URL"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image"`
	Phone string `json:"phone"`
}

type UserServiceResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    UserInfo `json:"data"`
}

func (c *UserServiceClient) GetUserByID(userID string) (*UserInfo, error) {
	url := fmt.Sprintf("%s/api/user/public/%s", c.baseURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get user failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response UserServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Data, nil
}

type UpdateRatingRequest struct {
	SellerID  string  `json:"seller_id"`
	Star      float64 `json:"star"`
	Operation string  `json:"operation"`
	OldStar   float64 `json:"old_star,omitempty"`
}

type UpdateRatingResponse struct {
	RatingCount   int     `json:"rating_count"`
	RatingAverage float64 `json:"rating_average"`
}

type RatingServiceResponse struct {
	Status  int                  `json:"status"`
	Message string               `json:"message"`
	Data    UpdateRatingResponse `json:"data"`
}

func (c *UserServiceClient) UpdateSellerRating(sellerID string, star float64, operation string, oldStar float64) error {
	url := fmt.Sprintf("%s/api/user/seller/rating", c.baseURL)

	requestBody := UpdateRatingRequest{
		SellerID:  sellerID,
		Star:      star,
		Operation: operation,
		OldStar:   oldStar,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update seller rating failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
