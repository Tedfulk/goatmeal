package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/database"
	"github.com/tedfulk/goatmeal/ui/theme"
	"github.com/tedfulk/goatmeal/utils/models"
)

type ShowMenuMsg struct{}

type ConversationItem struct {
	id        string
	title     string
	provider  string
	model     string
	messages  []database.Message
}

func (i ConversationItem) Title() string       { return i.title }
func (i ConversationItem) Description() string { return "" }
func (i ConversationItem) FilterValue() string { return i.title }

type KeyMap struct {
	Back key.Binding
	Delete key.Binding
	SwitchFocus key.Binding
	Export key.Binding
}

var DefaultKeyMap = KeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to menu"),
	),
	Delete: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "delete conversation"),
	),
	SwitchFocus: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch focus"),
	),
	Export: key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "export conversation"),
	),
}

type ConversationListView struct {
	db       *database.DB
	config   *config.Config
	list     list.Model
	messages []database.Message
	width    int
	height   int
	selected int
	keys     KeyMap
	viewport viewport.Model
	focused  string
}

type ResetTitleMsg struct{}

func NewConversationListView(db *database.DB, cfg *config.Config) *ConversationListView {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Conversations"
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = theme.BaseStyle.Title.
		Foreground(theme.CurrentTheme.Primary.GetColor()).
		Align(lipgloss.Center).
		Width(30)

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			DefaultKeyMap.Back,
			DefaultKeyMap.Delete,
			DefaultKeyMap.SwitchFocus,
			DefaultKeyMap.Export,
		}
	}
	l.AdditionalFullHelpKeys = l.AdditionalShortHelpKeys

	vp := viewport.New(104, 34)
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.CurrentTheme.Primary.GetColor()).
		Padding(1, 1)
	vp.MouseWheelEnabled = true

	view := &ConversationListView{
		db:       db,
		config:   cfg,
		list:     l,
		keys:     DefaultKeyMap,
		viewport: vp,
		focused:  "list",
	}

	// Load all conversations
	view.loadConversations()

	// Load messages for the first conversation if any exist
	if len(view.list.Items()) > 0 {
		firstConv := view.list.Items()[0].(ConversationItem)
		view.loadMessages(firstConv.id)
	}

	return view
}

func (c *ConversationListView) loadConversations() {
	// Load all conversations at once
	conversations, err := c.db.GetConversations(0, -1)
	if err != nil {
		return
	}

	var items []list.Item
	for _, conv := range conversations {
		items = append(items, ConversationItem{
			id:       conv.ID,
			title:    conv.Title,
			provider: conv.Provider,
			model:    conv.Model,
		})
	}
	c.list.SetItems(items)

	// If there are conversations, load messages for the first one
	if len(items) > 0 {
		firstConv := items[0].(ConversationItem)
		c.loadMessages(firstConv.id)
	}
}

func (c *ConversationListView) loadMessages(conversationID string) {
	messages, err := c.db.GetConversationMessages(conversationID)
	if err != nil {
		return
	}
	c.messages = messages

	// Get the conversation details to access the model name
	conversations, err := c.db.GetConversations(0, -1)
	if err != nil {
		return
	}
	
	var currentConv *database.Conversation
	for _, conv := range conversations {
		if conv.ID == conversationID {
			currentConv = &conv
			break
		}
	}

	// Update viewport content
	var content string
	if len(c.messages) > 0 {
		for _, msg := range c.messages {
			var prefix string
			if msg.Role == "user" {
				prefix = lipgloss.NewStyle().
					Foreground(theme.CurrentTheme.Message.UserText.GetColor()).
						Render(c.config.Settings.Username)
			} else if msg.Role == "search" {
				prefix = lipgloss.NewStyle().
					Foreground(theme.CurrentTheme.Message.AIText.GetColor()).
						Render("Tavily")
			} else {
				// Color model name with AIText color
				modelName := "AI"
				if currentConv.Provider == "tavily" {
					modelName = "Tavily"
				} else {
					modelName = models.StripModelsPrefix(currentConv.Model)
				}
				prefix = lipgloss.NewStyle().
					Foreground(theme.CurrentTheme.Message.AIText.GetColor()).
						Render(modelName)
			}

			// Render message content with Glamour if enabled
			msgContent := msg.Content
			if c.config.Settings.OutputGlamour && (msg.Role == "assistant" || msg.Role == "search") {
				if rendered, err := glamour.Render(msg.Content, "dark"); err == nil {
					msgContent = rendered
				}
			}

			content += prefix + "\n" + msgContent + "\n\n"
		}
	} else {
		content = "Select a conversation to view messages"
	}
	c.viewport.SetContent(content)
}

func (c *ConversationListView) exportConversation(id string) (tea.Cmd, error) {
	// Set the title immediately
	c.list.Title = "Exporting"
	
	// Create a command to reset the title after 1 second
	resetCmd := tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return ResetTitleMsg{}
	})

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %w", err)
	}

	// Get the conversation data
	conv, err := c.db.ExportConversation(id)
	if err != nil {
		return nil, fmt.Errorf("error exporting conversation: %w", err)
	}

	// Create the export data structure
	exportData := struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Provider  string    `json:"provider"`
		Model     string    `json:"model"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Messages  []struct {
			Role      string    `json:"role"`
			Content   string    `json:"content"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"messages"`
	}{
		ID:        conv.ID,
		Title:     conv.Title,
		Provider:  conv.Provider,
		Model:     conv.Model,
		CreatedAt: conv.CreatedAt,
		UpdatedAt: conv.UpdatedAt,
		Messages:  make([]struct {
			Role      string    `json:"role"`
			Content   string    `json:"content"`
			Timestamp time.Time `json:"timestamp"`
		}, len(conv.Messages)),
	}

	// Convert messages
	for i, msg := range conv.Messages {
		exportData.Messages[i] = struct {
			Role      string    `json:"role"`
			Content   string    `json:"content"`
			Timestamp time.Time `json:"timestamp"`
		}{
			Role:      msg.Role,
			Content:   msg.Content,
			Timestamp: msg.CreatedAt,
		}
	}

	// Create JSON data
	jsonData, err := json.MarshalIndent(exportData, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %w", err)
	}

	// Create filename using conversation title (sanitized) and timestamp
	sanitizedTitle := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' || r == '_' {
			return r
		}
		return '-'
	}, conv.Title)
	
	filename := fmt.Sprintf("%s/Downloads/conversation-%s-%s.json",
		homeDir,
		sanitizedTitle,
		time.Now().Format("2006-01-02-150405"))

	// Write the file
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return nil, fmt.Errorf("error writing file: %w", err)
	}

	return resetCmd, nil
}

func (c *ConversationListView) Update(msg tea.Msg) (*ConversationListView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ResetTitleMsg:
		c.list.Title = "Conversations"
		return c, nil

	case tea.MouseMsg:
		if c.focused == "messages" {
			if msg.Action == tea.MouseActionPress {
				if msg.Button == tea.MouseButtonWheelUp {
					c.viewport.LineUp(1)
				} else if msg.Button == tea.MouseButtonWheelDown {
					c.viewport.LineDown(1)
				}
			}
		}

	case tea.KeyMsg:
		// Handle tab key for focus switching
		if msg.String() == "tab" {
			if c.focused == "list" {
				c.focused = "messages"
			} else {
				c.focused = "list"
			}
			return c, nil
		}

		// Handle our key bindings first
		if key.Matches(msg, c.keys.Back) {
			return c, func() tea.Msg {
				return SetViewMsg{view: "chat"}
			}
		}

		if key.Matches(msg, c.keys.Delete) {
			if len(c.list.Items()) > 0 {
				selected := c.list.SelectedItem().(ConversationItem)
				if err := c.db.DeleteConversation(selected.id); err == nil {
					c.loadConversations()
					c.messages = nil
					c.viewport.SetContent("Select a conversation to view messages")
				}
			}
		}

		if key.Matches(msg, c.keys.Export) {
			if len(c.list.Items()) > 0 {
				selected := c.list.SelectedItem().(ConversationItem)
				if cmd, err := c.exportConversation(selected.id); err != nil {
					fmt.Printf("Error exporting conversation: %v\n", err)
				} else {
					return c, cmd
				}
			}
			return c, nil
		}

		// Only pass key events to the focused component
		if c.focused == "list" {
			var listCmd tea.Cmd
			c.list, listCmd = c.list.Update(msg)
			cmds = append(cmds, listCmd)
		} else {
			var vpCmd tea.Cmd
			c.viewport, vpCmd = c.viewport.Update(msg)
			cmds = append(cmds, vpCmd)
		}
	}

	// If the selected conversation changed, load its messages
	newSelected := c.list.Index()
	if newSelected != c.selected && len(c.list.Items()) > 0 {
		selected := c.list.SelectedItem().(ConversationItem)
		c.loadMessages(selected.id)
	}
	c.selected = newSelected

	return c, tea.Batch(cmds...)
}

func (c ConversationListView) View() string {
	// Left container (list)
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.CurrentTheme.Primary.GetColor()).
		Padding(1, 1).
		Width(34).
		Height(34)

	// Highlight the focused container's border
	if c.focused == "list" {
		listStyle = listStyle.BorderForeground(theme.CurrentTheme.Border.Active.GetColor())
	}

	// Style the viewport based on focus
	vpStyle := c.viewport.Style
	vpStyle = vpStyle.BorderForeground(theme.CurrentTheme.Primary.GetColor())
	if c.focused == "messages" {
		vpStyle = vpStyle.BorderForeground(theme.CurrentTheme.Border.Active.GetColor())
	}
	c.viewport.Style = vpStyle

	containers := lipgloss.JoinHorizontal(
		lipgloss.Left,
		listStyle.Render(c.list.View()),
		c.viewport.View(),
	)

	return lipgloss.Place(
		c.width,
		c.height,
		lipgloss.Left,
		lipgloss.Center,
		containers,
	)
}

func (c *ConversationListView) SetSize(width, height int) {
	c.width = width
	c.height = height
	c.list.SetSize(30, height - 4)
	
	c.viewport.Width = width - 37
	c.viewport.Height = height
} 
