package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tedfulk/goatmeal/services/providers"
)

const (
	baseURL          = "https://api.anthropic.com/v1"
	anthropicVersion = "2023-06-01"
)

// Provider implements the providers.Provider interface for Anthropic
type Provider struct {
	providers.BaseProvider
	client *http.Client
}

// NewProvider creates a new Anthropic provider
func NewProvider(apiKey string) *Provider {
	return &Provider{
		BaseProvider: providers.NewBaseProvider("anthropic", apiKey),
		client:      &http.Client{},
	}
}

// SendMessage sends a message to Anthropic and returns the response
func (p *Provider) SendMessage(ctx context.Context, message, systemPrompt, model string) (string, error) {
	url := fmt.Sprintf("%s/messages", baseURL)
	
	// Combine system prompt and user message
	messages := []map[string]string{
		{"role": "user", "content": fmt.Sprintf("%s\n\n%s", systemPrompt, message)},
	}

	payload := map[string]interface{}{
		"model":      model,
		"max_tokens": 8192,
		"messages":   messages,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.GetAPIKey())
	req.Header.Set("anthropic-version", anthropicVersion)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("no response from anthropic")
	}

	return result.Content[0].Text, nil
}

// ListModels returns a list of available Anthropic models
func (p *Provider) ListModels(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/models", baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("x-api-key", p.GetAPIKey())
	req.Header.Set("anthropic-version", anthropicVersion)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Data []struct {
			ID          string `json:"id"`
			Type        string `json:"type"`
			DisplayName string `json:"display_name"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	models := make([]string, 0, len(result.Data))
	for _, model := range result.Data {
		if model.Type == "model" {
			models = append(models, model.ID)
		}
	}

	return models, nil
}

// ValidateAPIKey checks if the API key is valid by attempting to list models
func (p *Provider) ValidateAPIKey(ctx context.Context) error {
	_, err := p.ListModels(ctx)
	return err
} 