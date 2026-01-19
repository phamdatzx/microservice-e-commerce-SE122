package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type NotificationServiceClient struct {
	baseURL string
	client  *http.Client
}

func NewNotificationServiceClient(baseURL string) *NotificationServiceClient {
	return &NotificationServiceClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CreateNotificationRequest represents the request payload for creating a notification
type CreateNotificationRequest struct {
	UserID  string                 `json:"userId"`
	Type    string                 `json:"type"`
	Title   string                 `json:"title"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// CreateNotification sends a notification to the notification service
func (c *NotificationServiceClient) CreateNotification(request CreateNotificationRequest) error {
	url := fmt.Sprintf("%s/api/notification", c.baseURL)

	bodyBytes, err := json.Marshal(request)
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create notification failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
