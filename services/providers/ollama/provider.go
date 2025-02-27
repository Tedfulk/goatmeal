// Package ollama provides integration with the Ollama API for local language model inference
package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tedfulk/goatmeal/services/providers"
)

const (
	// defaultBaseURL is the default URL for the Ollama API
	defaultBaseURL = "http://localhost:11434/api"
)

// Provider implements the providers.Provider interface for Ollama
type Provider struct {
	providers.OpenAICompatibleProvider
}

// NewProvider creates a new Ollama provider
// Although apiKey is not used by Ollama, it's included for interface compatibility
func NewProvider(apiKey string) *Provider {
	cfg := providers.OpenAICompatibleConfig{
		Name:    "ollama",
		APIKey:  "ollama",
		ModelFilter: func(modelID string) bool {
			return true
		},
	}
	return &Provider{
		OpenAICompatibleProvider: providers.NewOpenAICompatibleProvider(cfg),
	}
}

// ollamaModelsResponse represents the response structure from Ollama's /tags endpoint
type ollamaModelsResponse struct {
	Models []struct {
		Name string `json:"name"`
	} `json:"models"`
}

// ollamaChatResponse represents the response structure from Ollama's /chat endpoint
type ollamaChatResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

// ListModels returns a list of available Ollama models
func (p *Provider) ListModels(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/tags", defaultBaseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result ollamaModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	models := make([]string, 0, len(result.Models))
	for _, model := range result.Models {
		models = append(models, model.Name)
	}

	return models, nil
}

// SendMessage sends a chat message to Ollama and returns the response
func (p *Provider) SendMessage(ctx context.Context, message, systemPrompt, model string) (string, error) {
	url := fmt.Sprintf("%s/chat", defaultBaseURL)

	payload := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": message,
			},
		},
		"stream": false,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return result.Message.Content, nil
}
