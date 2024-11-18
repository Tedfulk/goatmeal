package ui

import (
	"fmt"
	"goatmeal/api"
	"goatmeal/config"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type ChatModel struct {
	viewport    viewport.Model
	messages    []api.Message
	style      Style
	width      int
	height     int
	ready      bool
	renderer   *glamour.TermRenderer
}

type Style struct {
	UserBubble      lipgloss.Style
	AssistantBubble lipgloss.Style
	Timestamp       lipgloss.Style
}

func NewChat(config *config.Config) (ChatModel, error) {
	colors := config.GetThemeColors()
	
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return ChatModel{}, fmt.Errorf("failed to initialize markdown renderer: %w", err)
	}

	return ChatModel{
		messages: make([]api.Message, 0),
		style:   Style{
			UserBubble: lipgloss.NewStyle().
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colors.UserBubble).
				Align(lipgloss.Right),
			
			AssistantBubble: lipgloss.NewStyle().
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colors.AssistantBubble).
				Align(lipgloss.Left),
			
			Timestamp: lipgloss.NewStyle().
				Foreground(colors.Timestamp).
				Align(lipgloss.Center),
		},
		renderer: renderer,
		width:    80,
	}, nil
}

func (m *ChatModel) AddMessage(msg api.Message) {
	m.messages = append(m.messages, msg)
	if m.ready {
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()
	}
}

func (m ChatModel) GetMessages() []api.Message {
	return m.messages
}

func (m ChatModel) Init() tea.Cmd {
	return nil
}

func (m ChatModel) Update(msg tea.Msg) (ChatModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+b", "pgup":
			m.viewport.HalfViewUp()
		case "ctrl+f", "pgdown":
			m.viewport.HalfViewDown()
		case "up", "k":
			m.viewport.LineUp(1)
		case "down", "j":
			m.viewport.LineDown(1)
		case "home":
			m.viewport.GotoTop()
		case "end":
			m.viewport.GotoBottom()
		}

	case tea.WindowSizeMsg:
		headerHeight := 0
		footerHeight := 6
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.renderMessages())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		m.width = msg.Width
		m.height = msg.Height
	}

	// Important: Update viewport and capture command
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ChatModel) View() string {
	if !m.ready {
		return "\nInitializing..."
	}

	return m.viewport.View()
}

func (m ChatModel) renderMessage(msg api.Message) string {
	timestamp := m.style.Timestamp.Render(msg.Timestamp.Format("15:04"))
	
	var content string
	var err error

	// Calculate max width for messages (70% of viewport width)
	maxWidth := m.width * 7 / 10

	if msg.Role == "assistant" {
		// Assistant messages on the left
		content, err = m.renderer.Render(msg.Content)
		if err != nil {
			content = fmt.Sprintf("Error rendering markdown: %v\nOriginal message:\n%s", err, msg.Content)
		}
		
		// Create left-aligned message block with padding
		messageBlock := lipgloss.NewStyle().
			Width(m.width).
			PaddingLeft(2).  // Add left padding for assistant messages
			Align(lipgloss.Left).  // Ensure left alignment
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					(m.style.AssistantBubble).
						MaxWidth(maxWidth).
						Render(content),
					lipgloss.NewStyle().
						PaddingLeft(4).  // Indent timestamp
						Render(timestamp),
				),
			)
		
		return messageBlock + "\n"

	} else {
		// User messages on the right
		// Create right-aligned message block
		messageBlock := lipgloss.NewStyle().
			Width(m.width).
			PaddingRight(2).  // Add right padding for user messages
			Align(lipgloss.Right).  // Ensure right alignment
			Render(
				lipgloss.JoinVertical(
					lipgloss.Right,
					(m.style.UserBubble).
						MaxWidth(maxWidth).
						Align(lipgloss.Right).  // Right align the message content
						Render(msg.Content),
					lipgloss.NewStyle().
						PaddingRight(4).  // Indent timestamp
						Render(timestamp),
				),
			)
		
		return messageBlock + "\n"
	}
}

func (m ChatModel) renderMessages() string {
	var sb strings.Builder

	for _, msg := range m.messages {
		sb.WriteString(m.renderMessage(msg))
	}

	return sb.String()
}
