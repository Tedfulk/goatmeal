package deepseek

import (
	"strings"

	"github.com/tedfulk/goatmeal/services/providers"
)

// Provider implements the providers.Provider interface for Deepseek
type Provider struct {
	providers.OpenAICompatibleProvider
}

// NewProvider creates a new Deepseek provider
func NewProvider(apiKey string) *Provider {
	return &Provider{
		OpenAICompatibleProvider: providers.NewOpenAICompatibleProvider(providers.OpenAICompatibleConfig{
			Name:    "deepseek",
			APIKey:  apiKey,
			ModelFilter: func(modelID string) bool {
				return !strings.Contains(strings.ToLower(modelID), "whisper")
			},
			ExtraHeaders: map[string]string{
				"Accept": "application/json",
			},
			ExtraParams: map[string]interface{}{
				"stream": false,
			},
		}),
	}
}
