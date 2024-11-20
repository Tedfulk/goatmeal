package ui

import (
	"fmt"
	"goatmeal/api"
	"goatmeal/config"
	"goatmeal/db"
	"strings"

	"github.com/charmbracelet/bubbles/key"
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
	db         db.ChatDB
	currentID  string
	focused    bool
}

type Style struct {
	UserBubble      lipgloss.Style
	AssistantBubble lipgloss.Style
	Timestamp       lipgloss.Style
}

func NewStyle(colors config.ThemeColors) Style {
	return Style{
		UserBubble: lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colors.UserBubble)).
			Foreground(lipgloss.Color(colors.UserText)),

		AssistantBubble: lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colors.AssistantBubble)).
			Foreground(lipgloss.Color(colors.AssistantText)),

		Timestamp: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Timestamp)).
			Italic(true).
			Faint(true),
	}
}

func NewChat(config *config.Config, database db.ChatDB, conversationID string) (ChatModel, error) {
	colors := config.GetThemeColors()
	
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return ChatModel{}, fmt.Errorf("failed to initialize markdown renderer: %w", err)
	}

	// Initialize viewport with default size
	vp := viewport.New(80, 20)
	vp.KeyMap = viewport.KeyMap{
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+f"),
			key.WithHelp("pgdn/ctrl+f", "page down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+b"),
			key.WithHelp("pgup/ctrl+b", "page up"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("ctrl+u", "half page up"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "half page down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
	}

	// Initialize an empty chat model
	chat := ChatModel{
		messages:   []api.Message{}, // Initialize with empty messages
		style:     NewStyle(colors),
		renderer:  renderer,
		db:        database,
		currentID: conversationID,
		viewport:  vp,
		ready:     true,
	}

	// Only load messages if we have a valid conversation ID
	if conversationID != "" {
		dbMessages, err := database.GetMessages(conversationID)
		if err != nil {
			return ChatModel{}, fmt.Errorf("failed to load messages: %w", err)
		}

		// Convert db.Message to api.Message
		for _, msg := range dbMessages {
			messages := api.Message{
				Role:      msg.Role,
				Content:   msg.Content,
				Timestamp: msg.CreatedAt,
			}
			// Add each message to the chat
			if err := chat.AddMessage(messages); err != nil {
				return ChatModel{}, fmt.Errorf("failed to add message: %w", err)
			}
		}
	}

	// Set initial content (empty for new chat, or loaded messages for existing chat)
	chat.viewport.SetContent(chat.renderMessages())

	return chat, nil
}

func (m *ChatModel) AddMessage(msg api.Message) error {
	// Only try to add to database if we have a conversation ID
	if m.currentID != "" {
		// Add message to the database
		dbMsg := db.Message{
			ConversationID: m.currentID,
			Role:          msg.Role,
			Content:       msg.Content,
			CreatedAt:     msg.Timestamp,
		}

		if err := m.db.AddMessage(m.currentID, dbMsg); err != nil {
			return fmt.Errorf("failed to persist message: %w", err)
		}
	}

	// Add to in-memory messages regardless of database state
	m.messages = append(m.messages, msg)
	
	// Update viewport content
	m.viewport.SetContent(m.renderMessages())
	m.viewport.GotoBottom() // Make sure we scroll to the latest message

	return nil
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
		// Handle quitting the program
		if msg.String() == "q" {
			return m, tea.Quit
		}

		// Handle new chat shortcut
		if msg.String() == "shift+:" {
			return m, func() tea.Msg { return NewChatMsg{} }
		}

		// Handle focus switching
		if msg.String() == "tab" {
			m.focused = !m.focused
			return m, nil
		}

		// If viewport is focused, handle scrolling
		if m.focused {
			switch msg.String() {
			case "up", "k":
				m.viewport.LineUp(1)
			case "down", "j":
				m.viewport.LineDown(1)
			case "ctrl+b", "pgup":
				m.viewport.HalfViewUp()
			case "ctrl+f", "pgdown":
				m.viewport.HalfViewDown()
			case "home":
				m.viewport.GotoTop()
			case "end":
				m.viewport.GotoBottom()
			}
		}

	case tea.WindowSizeMsg:
		headerHeight := 1
		footerHeight := 3
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
		
		// Re-render messages with new width
		m.viewport.SetContent(m.renderMessages())
	}

	// Always update viewport regardless of focus state
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ChatModel) View() string {
	if !m.ready {
		return "\nInitializing..."
	}

	style := lipgloss.NewStyle()

	return style.Render(m.viewport.View())
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
