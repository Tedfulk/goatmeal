package search

import "strings"

// StripPrefix removes search-related prefixes from content
func StripPrefix(content string) string {
    prefixes := []string{
        "🔍 Searching for: ",
        "🔍+ Enhanced search: ",
    }
    
    for _, prefix := range prefixes {
        if strings.HasPrefix(content, prefix) {
            return strings.TrimPrefix(content, prefix)
        }
    }
    return content
} 