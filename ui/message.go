package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
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

// View renders the message with its ID and content
func (m Message) View(width int) string {
	// Format timestamp with username/model
	var prefix string
	if m.Type == UserMessage {
		prefix = m.Config.Settings.Username
	} else {
		prefix = m.Config.CurrentModel
	}
	
	timestampStr := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render(fmt.Sprintf("%s • #%d • %s", prefix, m.ID, m.Timestamp.Format("15:04")))

	baseStyle := lipgloss.NewStyle().
		Padding(1).
		BorderStyle(lipgloss.RoundedBorder())

	if m.Type == UserMessage {
		contentStyle := lipgloss.NewStyle().Foreground(primaryColor)
		
		renderedContent := contentStyle.Render(m.Content)
		
		messageBox := baseStyle.
			BorderForeground(primaryColor).
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
		contentStyle := lipgloss.NewStyle().
			Width(width - 12).
			Align(lipgloss.Left).
			Foreground(secondaryColor)

		content := baseStyle.
			BorderForeground(secondaryColor).
			Render(contentStyle.Render(m.Content))

		// Create a container for message and timestamp
		return lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			timestampStr,
		)
	}
} 