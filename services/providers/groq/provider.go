package groq

import (
	"strings"

	"github.com/tedfulk/goatmeal/services/providers"
)

// Provider implements the providers.Provider interface for Groq
type Provider struct {
	providers.OpenAICompatibleProvider
}

// NewProvider creates a new Groq provider
func NewProvider(apiKey string) *Provider {
	return &Provider{
		OpenAICompatibleProvider: providers.NewOpenAICompatibleProvider(providers.OpenAICompatibleConfig{
			Name:    "groq",
			APIKey:  apiKey,
			ModelFilter: func(modelID string) bool {
				return !strings.Contains(strings.ToLower(modelID), "whisper")
			},
		}),
	}
} 