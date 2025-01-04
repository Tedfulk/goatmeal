package providers

import "context"

// Provider defines the interface that all AI chat providers must implement
type Provider interface {
	// SendMessage sends a message to the AI provider and returns the response
	SendMessage(ctx context.Context, message, systemPrompt, model string) (string, error)
	
	// ListModels returns a list of available models for this provider
	ListModels(ctx context.Context) ([]string, error)
	
	// GetName returns the name of the provider
	GetName() string
	
	// ValidateAPIKey checks if the provided API key is valid
	ValidateAPIKey(ctx context.Context) error
}

// BaseProvider implements common functionality for all providers
type BaseProvider struct {
	name    string
	apiKey  string
}

// NewBaseProvider creates a new BaseProvider instance
func NewBaseProvider(name, apiKey string) BaseProvider {
	return BaseProvider{
		name:    name,
		apiKey:  apiKey,
	}
}

// GetName returns the provider name
func (p BaseProvider) GetName() string {
	return p.name
}

// GetAPIKey returns the API key
func (p BaseProvider) GetAPIKey() string {
	return p.apiKey
} 