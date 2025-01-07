package gemini

import (
	"context"
	"fmt"
	"mime"
	"path/filepath"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/tedfulk/goatmeal/services/providers"
	"google.golang.org/api/option"
)

// Provider implements the providers.Provider interface for Gemini
type Provider struct {
	providers.BaseProvider
	client *genai.Client
}

// NewProvider creates a new Gemini provider
func NewProvider(apiKey string) *Provider {
	return &Provider{
		BaseProvider: providers.NewBaseProvider("gemini", apiKey),
	}
}

// SendMessage sends a message to Gemini and returns the response
func (p *Provider) SendMessage(ctx context.Context, message, systemPrompt, model string) (string, error) {
	if p.client == nil {
		client, err := genai.NewClient(ctx, option.WithAPIKey(p.GetAPIKey()))
		if err != nil {
			return "", fmt.Errorf("error creating Gemini client: %w", err)
		}
		p.client = client
	}

	// Create a new model instance
	m := p.client.GenerativeModel(model)

	// Set system prompt if provided
	if systemPrompt != "" {
		m.SystemInstruction = genai.NewUserContent(genai.Text(systemPrompt))
	}

	// Configure the model
	m.SetTemperature(0.2)
	m.SetTopK(40)
	m.SetTopP(0.95)
	m.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
	}

	// Create chat session
	cs := m.StartChat()

	// Send the message
	resp, err := cs.SendMessage(ctx, genai.Text(message))
	if err != nil {
		return "", fmt.Errorf("error sending message: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	// Get the response text
	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", fmt.Errorf("unexpected response type from Gemini")
	}

	return string(text), nil
}

// SendMessageWithImage sends a message with an image to Gemini and returns the response
func (p *Provider) SendMessageWithImage(ctx context.Context, message string, imageData []byte, imagePath string, systemPrompt, model string) (string, error) {
	if p.client == nil {
		client, err := genai.NewClient(ctx, option.WithAPIKey(p.GetAPIKey()))
		if err != nil {
			return "", fmt.Errorf("error creating Gemini client: %w", err)
		}
		p.client = client
	}

	// Create a new model instance
	m := p.client.GenerativeModel(model)

	// Set system prompt if provided
	if systemPrompt != "" {
		m.SystemInstruction = genai.NewUserContent(genai.Text(systemPrompt))
	}

	// Configure the model
	m.SetTemperature(0.2)
	m.SetTopK(40)
	m.SetTopP(0.95)
	m.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
	}

	// Determine MIME type from file extension
	ext := strings.ToLower(filepath.Ext(imagePath))
	if ext == "" {
		ext = ".jpg" // default to jpg if no extension
	}
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "image/jpeg" // default to jpeg if no MIME type found
	}

	// Create parts for the message
	parts := []genai.Part{
		genai.ImageData(mimeType, imageData),
		genai.Text(message),
	}

	// Generate content with both image and text
	resp, err := m.GenerateContent(ctx, parts...)
	if err != nil {
		return "", fmt.Errorf("error generating content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	// Get the response text
	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", fmt.Errorf("unexpected response type from Gemini")
	}

	return string(text), nil
}

// ListModels returns a list of available Gemini models
func (p *Provider) ListModels(ctx context.Context) ([]string, error) {
	if p.client == nil {
		client, err := genai.NewClient(ctx, option.WithAPIKey(p.GetAPIKey()))
		if err != nil {
			return nil, fmt.Errorf("error creating Gemini client: %w", err)
		}
		p.client = client
	}

	// Get all available models
	models := make([]string, 0)
	iter := p.client.ListModels(ctx)
	for {
		model, err := iter.Next()
		if err != nil {
			break
		}
		// Only include Gemini models
		if model.Name != "" && model.Name != "models/gemini-pro-vision" {
			models = append(models, model.Name)
		}
	}

	if len(models) == 0 {
		// Fallback to known models if list fails
		return []string{
			"gemini-pro",
			"gemini-pro-vision",
		}, nil
	}

	return models, nil
}

// ValidateAPIKey checks if the API key is valid by attempting to create a client
func (p *Provider) ValidateAPIKey(ctx context.Context) error {
	client, err := genai.NewClient(ctx, option.WithAPIKey(p.GetAPIKey()))
	if err != nil {
		return fmt.Errorf("invalid API key: %w", err)
	}
	defer client.Close()
	return nil
} 