package ui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/ui/theme"
	"github.com/tedfulk/goatmeal/utils/models"
)

// MessageType represents the type of message (user or provider)
type MessageType int

const (
	UserMessage MessageType = iota
	ProviderMessage
	SearchMessage
)

// Message represents a single message in the conversation
type Message struct {
	ID        int
	Type      MessageType
	Content   string
	Timestamp time.Time
	Config    *config.Config
	codeBlocks []CodeBlock  // Store the code blocks when message is created
}

// NewMessage creates a new message
func NewMessage(id int, msgType MessageType, content string, cfg *config.Config, getNextBlockNum func() int) Message {
	msg := Message{
		ID:        id,
		Type:      msgType,
		Content:   content,
		Timestamp: time.Now(),
		Config:    cfg,
	}
	
	// If it's a provider message, process code blocks immediately
	if msgType == ProviderMessage {
		msg.codeBlocks = msg.processCodeBlocks(getNextBlockNum)
	}
	
	return msg
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
	} else if m.Type == SearchMessage {
		prefix = "Tavily"
	} else {
		prefix = models.StripModelsPrefix(m.Config.CurrentModel)
	}
	
	timestampStyle := lipgloss.NewStyle().
		Foreground(theme.CurrentTheme.Message.Timestamp.GetColor())
	
	timestampStr := timestampStyle.Render(
		fmt.Sprintf("%s • #m%d • %s", prefix, m.ID, m.Timestamp.Format("15:04")),
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
		var renderedContent string
		if m.Config.Settings.OutputGlamour {
			content := m.Content
			
			// Use stored code blocks to add numbers
			re := regexp.MustCompile("(?ms)```(.+?)```")
			blockIndex := 0
			content = re.ReplaceAllStringFunc(content, func(block string) string {
				if blockIndex < len(m.codeBlocks) {
					blockNum := m.codeBlocks[blockIndex].Number
					blockIndex++
					return block + fmt.Sprintf("\n[#b%d]\n", blockNum)
				}
				return block
			})

			glamourStyle := "dark"
			renderer, err := glamour.NewTermRenderer(
				glamour.WithStylePath(glamourStyle),
				glamour.WithWordWrap(120),
			)
			if err == nil {
				if rendered, err := renderer.Render(content); err == nil {
					renderedContent = rendered
				} else {
					renderedContent = content
				}
			} else {
				renderedContent = content
			}
		} else {
			// Similar changes for non-glamour rendering
			content := m.Content
			blockIndex := 0
			
			content = regexp.MustCompile("(?ms)```(.+?)```").ReplaceAllStringFunc(content, func(block string) string {
				if blockIndex < len(m.codeBlocks) {
					blockNum := m.codeBlocks[blockIndex].Number
					blockIndex++
					return block + fmt.Sprintf("\n[#b%d]\n", blockNum)
				}
				return block
			})
			
			renderedContent = content
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

// Add this struct after the Message struct
type CodeBlock struct {
	Number  int
	Content string
}

// Add method to process code blocks
func (m *Message) processCodeBlocks(getNextBlockNum func() int) []CodeBlock {
	var blocks []CodeBlock
	
	re := regexp.MustCompile("(?ms)```(.+?)```")
	matches := re.FindAllStringSubmatch(m.Content, -1)
	
	for _, match := range matches {
		blockNum := getNextBlockNum()
		content := strings.TrimSpace(match[1])
		if idx := strings.Index(content, "\n"); idx != -1 {
			content = content[idx+1:]
		}
		blocks = append(blocks, CodeBlock{
			Number:  blockNum,
			Content: content,
		})
	}
	
	return blocks
}

// Update ExtractCodeBlocks to use stored blocks
func (m Message) ExtractCodeBlocks() []CodeBlock {
	return m.codeBlocks
} 