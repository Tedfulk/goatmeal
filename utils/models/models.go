package models

import "strings"

// StripModelsPrefix removes the "models/" prefix from model names if present
func StripModelsPrefix(model string) string {
	if strings.HasPrefix(model, "models/") {
		return strings.TrimPrefix(model, "models/")
	}
	return model
} 