package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/dto"
	"os"
	"time"
)

type UserServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserServiceClient() *UserServiceClient {
	baseURL := os.Getenv("USER_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080" // Default for local development
	}

	return &UserServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *UserServiceClient) GetUserByID(userID string) (*dto.UserResponse, error) {
	url := fmt.Sprintf("%s/api/user/public/%s", c.baseURL, userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user-service returned status %d", resp.StatusCode)
	}

	var response struct {
		Data dto.UserResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Data, nil
}
