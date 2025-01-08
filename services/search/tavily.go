package search


import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const tavilyAPIEndpoint = "https://api.tavily.com/search"

type SearchRequest struct {
	APIKey         string   `json:"api_key"`
	Query          string   `json:"query"`
	IncludeAnswer  bool     `json:"include_answer"`
	IncludeDomains []string `json:"include_domains,omitempty"`
}

type SearchResponse struct {
	Query   string         `json:"query"`
	Answer  interface{}    `json:"answer"`
	Results []SearchResult `json:"results"`
	ResponseTime float64   `json:"response_time"`
}

type SearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

type TavilyClient struct {
	apiKey string
}

func NewClient(apiKey string) *TavilyClient {
	return &TavilyClient{apiKey: apiKey}
}

func (c *TavilyClient) Search(query string, domains []string) (*SearchResponse, error) {
	reqBody := SearchRequest{
		APIKey:         c.apiKey,
		Query:          query,
		IncludeAnswer:  true,
		IncludeDomains: domains,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", tavilyAPIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &searchResp, nil
} 