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

type AIServiceClient struct {
	baseURL string
	client  *http.Client
}

func NewAIServiceClient() *AIServiceClient {
	return &AIServiceClient{
		baseURL: os.Getenv("AI_SERVICE_URL"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type aiRecommendationRequest struct {
	UserID string `json:"user_id"`
	Limit  int    `json:"limit"`
}

type aiRecommendedItem struct {
	ID    string                 `json:"id"`
	Score float64                `json:"score"`
	Payload map[string]interface{} `json:"payload"`
}

type aiRecommendationResponse struct {
	UserID string              `json:"user_id"`
	Items  []aiRecommendedItem `json:"items"`
}

// GetUserRecommendations calls the ai-service and returns a list of recommended product IDs
// in the order returned by the AI (descending score).
func (c *AIServiceClient) GetUserRecommendations(userID string, limit int) ([]string, error) {
	url := fmt.Sprintf("%s/api/ai/users/recommendations", c.baseURL)

	reqBody := aiRecommendationRequest{
		UserID: userID,
		Limit:  limit,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ai recommendation request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create ai recommendation request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call ai-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ai-service returned status %d: %s", resp.StatusCode, string(body))
	}

	var aiResp aiRecommendationResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return nil, fmt.Errorf("failed to decode ai-service response: %w", err)
	}

	ids := make([]string, 0, len(aiResp.Items))
	for _, item := range aiResp.Items {
		ids = append(ids, item.ID)
	}

	return ids, nil
}
