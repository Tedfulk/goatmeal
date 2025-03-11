package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/tedfulk/goatmeal/services/providers"
	"github.com/tedfulk/goatmeal/utils/location"
	"github.com/tedfulk/goatmeal/utils/prompts"
)

type EnhanceType string

const (
	WebSearch     EnhanceType = "web"
	Programming   EnhanceType = "programming"
)

type QueryEnhancer struct {
    apiKey string
}

func NewQueryEnhancer(apiKey string) *QueryEnhancer {
    return &QueryEnhancer{
        apiKey: apiKey,
    }
}

func (qe *QueryEnhancer) Enhance(query string, enhanceType EnhanceType) (string, error) {
    if qe.apiKey == "" {
        return query, fmt.Errorf("groq API key not found")
    }
    
    cfg := providers.OpenAICompatibleConfig{
        Name:    "groq",
        APIKey:  qe.apiKey,
    }
    provider := providers.NewOpenAICompatibleProvider(cfg)
    
    var fullPrompt string
    switch enhanceType {
    case WebSearch:
        locationInfo := location.GetFormattedLocationAndTime()
        fullPrompt = prompts.GetEnhanceSearchPrompt(locationInfo, query)
    case Programming:
        fullPrompt = prompts.GetEnhanceProgrammingPrompt(query)
    default:
        return query, fmt.Errorf("invalid enhancement type")
    }
    
    enhancedQuery, err := provider.SendMessage(context.Background(), fullPrompt, "", "llama-3.3-70b-versatile")
    if err != nil {
        return query, fmt.Errorf("failed to enhance query: %v", err)
    }
    
    // Only extract the query for web searches, programming queries can keep their full response
    if enhanceType == WebSearch {
        extractPrompt := prompts.GetExtractQueryPrompt(enhancedQuery)
        cleanQuery, err := provider.SendMessage(context.Background(), extractPrompt, "", "llama-3.3-70b-versatile")
        if err != nil {
            return query, fmt.Errorf("failed to clean enhanced query: %v", err)
        }
        enhancedQuery = cleanQuery
    }
    
    return strings.TrimSpace(enhancedQuery), nil
} 