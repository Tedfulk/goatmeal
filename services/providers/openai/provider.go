package openai

import (
	"strings"

	"github.com/tedfulk/goatmeal/services/providers"
)

// Provider implements the providers.Provider interface for OpenAI
type Provider struct {
	providers.OpenAICompatibleProvider
}

// NewProvider creates a new OpenAI provider
func NewProvider(apiKey string) *Provider {
	return &Provider{
		OpenAICompatibleProvider: providers.NewOpenAICompatibleProvider(providers.OpenAICompatibleConfig{
			Name:    "openai",
			APIKey:  apiKey,
			ModelFilter: func(modelID string) bool {
				return strings.HasPrefix(modelID, "gpt")
			},
		}),
	}
} 