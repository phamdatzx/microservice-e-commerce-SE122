package client

import (
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
