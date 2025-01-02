package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// OpenAICompatibleProvider implements common functionality for providers with OpenAI-compatible APIs
type OpenAICompatibleProvider struct {
	BaseProvider
	baseURL     string
	client      *http.Client
	modelFilter func(string) bool
	extraHeaders map[string]string
	extraParams map[string]interface{}
}

// getBaseURL returns the base URL for a given provider
func getBaseURL(providerName string) string {
	baseURLs := map[string]string{
		"openai":   "https://api.openai.com/v1",
		"groq":     "https://api.groq.com/openai/v1",
		"deepseek": "https://api.deepseek.com",
	}
	
	if url, ok := baseURLs[providerName]; ok {
		return url
	}
	return baseURLs["openai"] // default to OpenAI if provider not found
}

// OpenAICompatibleConfig represents the configuration for an OpenAI-compatible provider
type OpenAICompatibleConfig struct {
	Name         string
	APIKey       string
	ModelFilter  func(string) bool
	ExtraHeaders map[string]string
	ExtraParams  map[string]interface{}
}

// NewOpenAICompatibleProvider creates a new OpenAI-compatible provider
func NewOpenAICompatibleProvider(cfg OpenAICompatibleConfig) OpenAICompatibleProvider {
	if cfg.ExtraHeaders == nil {
		cfg.ExtraHeaders = make(map[string]string)
	}
	if cfg.ExtraParams == nil {
		cfg.ExtraParams = make(map[string]interface{})
	}

	return OpenAICompatibleProvider{
		BaseProvider:  NewBaseProvider(cfg.Name, cfg.APIKey),
		baseURL:      getBaseURL(cfg.Name),
		client:       &http.Client{},
		modelFilter:  cfg.ModelFilter,
		extraHeaders: cfg.ExtraHeaders,
		extraParams:  cfg.ExtraParams,
	}
}

// SendMessage sends a message to the provider and returns the response
func (p OpenAICompatibleProvider) SendMessage(ctx context.Context, message, systemPrompt, model string) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", p.baseURL)
	
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
		"temperature": 0.0,
	}

	// Add any extra parameters
	for k, v := range p.extraParams {
		payload[k] = v
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.GetAPIKey()))
	
	// Add any extra headers
	for k, v := range p.extraHeaders {
		req.Header.Set(k, v)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from %s", p.GetName())
	}

	return result.Choices[0].Message.Content, nil
}

// ListModels returns a list of available models
func (p OpenAICompatibleProvider) ListModels(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/models", p.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.GetAPIKey()))
	
	// Add any extra headers
	for k, v := range p.extraHeaders {
		req.Header.Set(k, v)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Try OpenAI format first
	var openAIResult struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&openAIResult); err == nil && len(openAIResult.Data) > 0 {
		models := make([]string, 0, len(openAIResult.Data))
		for _, model := range openAIResult.Data {
			if p.modelFilter(model.ID) {
				models = append(models, model.ID)
			}
		}
		return models, nil
	}

	// Reset response body for next read
	resp.Body.Close()
	req, _ = http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.GetAPIKey()))
	for k, v := range p.extraHeaders {
		req.Header.Set(k, v)
	}
	resp, err = p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Try Deepseek format
	var deepseekResult struct {
		Models []struct {
			ID string `json:"id"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&deepseekResult); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	models := make([]string, 0, len(deepseekResult.Models))
	for _, model := range deepseekResult.Models {
		if p.modelFilter(model.ID) {
			models = append(models, model.ID)
		}
	}
	return models, nil
}

// ValidateAPIKey checks if the API key is valid
func (p OpenAICompatibleProvider) ValidateAPIKey(ctx context.Context) error {
	_, err := p.ListModels(ctx)
	return err
} 