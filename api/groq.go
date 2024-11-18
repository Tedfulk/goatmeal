package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goatmeal/config"
	"net/http"
	"time"
)

type GroqClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	model      string
	systemMsg  string
}

type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"-"` // For UI display, not sent to API
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func NewGroqClient(config *config.Config) (*GroqClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key not found in config")
	}

	return &GroqClient{
		apiKey:     config.APIKey,
		baseURL:    "https://api.groq.com/openai/v1/chat/completions",
		httpClient: &http.Client{},
		model:      config.DefaultModel,
		systemMsg:  config.SystemPrompt,
	}, nil
}

func (c *GroqClient) SendMessage(userMessage string, conversationHistory []Message) (*Message, error) {
	// Prepare messages with system prompt and history
	messages := []Message{
		{Role: "system", Content: c.systemMsg},
	}
	messages = append(messages, conversationHistory...)
	messages = append(messages, Message{
		Role:      "user",
		Content:   userMessage,
		Timestamp: time.Now(),
	})

	reqBody := ChatRequest{
		Model:    c.model,
		Messages: messages,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	var response ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Add timestamp to the response message
	response.Choices[0].Message.Timestamp = time.Now()
	return &response.Choices[0].Message, nil
}
