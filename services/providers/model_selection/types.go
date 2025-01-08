package model_selection

import (
	"context"
	"fmt"

	"github.com/tedfulk/goatmeal/services/providers"
	"github.com/tedfulk/goatmeal/services/providers/anthropic"
	"github.com/tedfulk/goatmeal/services/providers/deepseek"
	"github.com/tedfulk/goatmeal/services/providers/gemini"
	"github.com/tedfulk/goatmeal/services/providers/groq"
	"github.com/tedfulk/goatmeal/services/providers/openai"
)

// Model represents an AI model
type Model struct {
	ID string
}

func (m Model) Title() string       { return m.ID }
func (m Model) Description() string { return "" }
func (m Model) FilterValue() string { return m.ID }

// FetchModels fetches available models for the selected provider
func FetchModels(provider, apiKey string) ([]string, error) {
	var p providers.Provider
	switch provider {
	case "openai":
		p = openai.NewProvider(apiKey)
	case "groq":
		p = groq.NewProvider(apiKey)
	case "deepseek":
		p = deepseek.NewProvider(apiKey)
	case "anthropic":
		p = anthropic.NewProvider(apiKey)
	case "gemini":
		p = gemini.NewProvider(apiKey)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	return p.ListModels(context.Background())
} 