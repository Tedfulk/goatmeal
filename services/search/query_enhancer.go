package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/tedfulk/goatmeal/services/providers"
	"github.com/tedfulk/goatmeal/utils/location"
	"github.com/tedfulk/goatmeal/utils/prompts"
)

type QueryEnhancer struct {
    apiKey string
}

func NewQueryEnhancer(apiKey string) *QueryEnhancer {
    return &QueryEnhancer{
        apiKey: apiKey,
    }
}

func (qe *QueryEnhancer) Enhance(query string) (string, error) {
    locationInfo := location.GetFormattedLocationAndTime()
    
    fullPrompt := prompts.GetEnhanceSearchPrompt(locationInfo, query)
    
    if qe.apiKey == "" {
        return query, fmt.Errorf("groq API key not found")
    }
    
    cfg := providers.OpenAICompatibleConfig{
        Name:    "groq",
        APIKey:  qe.apiKey,
    }
    provider := providers.NewOpenAICompatibleProvider(cfg)
    
    enhancedQuery, err := provider.SendMessage(context.Background(), fullPrompt, "", "llama-3.3-70b-versatile")
    if err != nil {
        return query, fmt.Errorf("failed to enhance query: %v", err)
    }
    
    extractPrompt := prompts.GetExtractQueryPrompt(enhancedQuery)
    cleanQuery, err := provider.SendMessage(context.Background(), extractPrompt, "", "llama-3.3-70b-versatile")
    if err != nil {
        return query, fmt.Errorf("failed to clean enhanced query: %v", err)
    }
    
    return strings.TrimSpace(cleanQuery), nil
} 