package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/ui/theme"
)

// MessageType represents the type of message (user or provider)
type MessageType int

const (
	UserMessage MessageType = iota
	ProviderMessage
)

// Message represents a single message in the conversation
type Message struct {
	ID        int
	Type      MessageType
	Content   string
	Timestamp time.Time
	Config    *config.Config
}

// NewMessage creates a new message
func NewMessage(id int, msgType MessageType, content string, cfg *config.Config) Message {
	return Message{
		ID:        id,
		Type:      msgType,
		Content:   content,
		Timestamp: time.Now(),
		Config:    cfg,
	}
}

// wordWrap wraps text at the specified width
func wordWrap(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return strings.Join(lines, "\n")
}

// View renders the message with its ID and content
func (m Message) View(width int) string {
	// Format timestamp with username/model
	var prefix string
	if m.Type == UserMessage {
		prefix = m.Config.Settings.Username
	} else {
		prefix = m.Config.CurrentModel
	}
	
	timestampStyle := lipgloss.NewStyle().
		Foreground(theme.CurrentTheme.Message.Timestamp.GetColor())
	
	timestampStr := timestampStyle.Render(
		fmt.Sprintf("%s • #%d • %s", prefix, m.ID, m.Timestamp.Format("15:04")),
	)

	baseStyle := theme.BaseStyle.Message

	if m.Type == UserMessage {
		// Calculate available width for content
		contentWidth := width - 16
		
		wrappedContent := wordWrap(m.Content, contentWidth)
		
		contentStyle := lipgloss.NewStyle().
			Foreground(theme.CurrentTheme.Message.UserText.GetColor())
		
		renderedContent := contentStyle.Render(wrappedContent)
		
		messageBox := baseStyle.
			BorderForeground(theme.CurrentTheme.Primary.GetColor()).
			Render(renderedContent)

		return lipgloss.NewStyle().
			Width(width - 6).
			Align(lipgloss.Right).
			Render(lipgloss.JoinVertical(
				lipgloss.Right,
				messageBox,
				timestampStr,
			))
	} else {
		// AI message rendering with optional Glamour
		var renderedContent string
		if m.Config.Settings.OutputGlamour {
			// Use Glamour for markdown rendering
			glamourStyle := "dark"
			if rendered, err := glamour.Render(m.Content, glamourStyle); err == nil {
				renderedContent = rendered
			} else {
				// Fallback to plain text if Glamour fails
				renderedContent = m.Content
			}
		} else {
			renderedContent = m.Content
		}

		contentStyle := lipgloss.NewStyle().
			Width(width - 12).
			Align(lipgloss.Left).
			Foreground(theme.CurrentTheme.Message.AIText.GetColor())

		content := baseStyle.
			BorderForeground(theme.CurrentTheme.Secondary.GetColor()).
			Render(contentStyle.Render(renderedContent))

		return lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			timestampStr,
		)
	}
} 