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

type ChatBotClient struct {
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
}

// ChatMessage represents a single message in the conversation
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the request body for OpenRouter API
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
}

// ChatChoice represents a choice in the API response
type ChatChoice struct {
	Index   int         `json:"index"`
	Message ChatMessage `json:"message"`
}

// ChatResponse represents the response from OpenRouter API
type ChatResponse struct {
	ID      string       `json:"id"`
	Choices []ChatChoice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func NewChatBotClient() *ChatBotClient {
	return &ChatBotClient{
		baseURL: os.Getenv("OPENROUTER_BASE_URL"),
		model:   os.Getenv("OPENROUTER_MODEL"),
		apiKey:  os.Getenv("OPENROUTER_API_KEY"),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendMessage sends a chat request to OpenRouter API
func (c *ChatBotClient) SendMessage(userMessage string, productContext string) (string, error) {
	// Prepare the request body
	requestBody := ChatRequest{
		Model: c.model,
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: "Bạn là chatbot hỗ trợ khách hàng của một shop bán hàng online. Trả lời ngắn gọn, rõ ràng.",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Bạn là chatbot hỗ trợ khách hàng của sàn thương mại điện tử. Chỉ được dùng thông tin bên dưới, không được bịa.\n%s\n=== CÂU HỎI KHÁCH HÀNG ===\n%s", productContext, userMessage),
			},
		},
		MaxTokens:   200,
		Temperature: 0.7,
	}

	// Marshal request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Unmarshal response
	var chatResponse ChatResponse
	if err := json.Unmarshal(body, &chatResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract response message
	if len(chatResponse.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return chatResponse.Choices[0].Message.Content, nil
}
