package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tedfulk/goatmeal/config"
)

type GroqClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	model      string
	systemMsg  string
	logger     *log.Logger
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

	// Open or create the log file
	logFile, err := os.OpenFile("groq_api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %w", err)
	}

	// Create a logger
	logger := log.New(logFile, "GROQ_API ", log.LstdFlags)

	return &GroqClient{
		apiKey:     config.APIKey,
		baseURL:    "https://api.groq.com/openai/v1/chat/completions",
		httpClient: &http.Client{},
		model:      config.DefaultModel,
		systemMsg:  config.SystemPrompt,
		logger:     logger,
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

	// Log the request payload
	c.logger.Printf("Sending request to Groq API: %+v\n", reqBody)

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

	// Log the response payload
	c.logger.Printf("Received response from Groq API: %+v\n", response)

	// Add timestamp to the response message
	response.Choices[0].Message.Timestamp = time.Now()
	return &response.Choices[0].Message, nil
}
