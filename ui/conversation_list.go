package ui

import (
	"fmt"
	"goatmeal/db"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define key mappings for the conversation list
type conversationListKeyMap struct {
    Up     key.Binding
    Down   key.Binding
    Select key.Binding
    Back   key.Binding
    Quit   key.Binding
}

var conversationListKeys = conversationListKeyMap{
    Up: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("↑/k", "up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("down", "j"),
        key.WithHelp("↓/j", "down"),
    ),
    Select: key.NewBinding(
        key.WithKeys("enter"),
        key.WithHelp("enter", "select conversation"),
    ),
    Back: key.NewBinding(
        key.WithKeys("esc", "tab"),
        key.WithHelp("esc/tab", "back"),
    ),
    Quit: key.NewBinding(
        key.WithKeys("q", "ctrl+c"),
        key.WithHelp("q/ctrl+c", "quit"),
    ),
}

// Item represents a conversation in the list
type conversationItem struct {
    id        string
    title     string
    createdAt time.Time
}

// Implement list.Item interface
func (i conversationItem) Title() string       { return i.title }
func (i conversationItem) Description() string { 
    return fmt.Sprintf("Created: %s", i.createdAt.Format("2006-01-02 15:04:05"))
}
func (i conversationItem) FilterValue() string { return i.title }

type ConversationListModel struct {
    list     list.Model
    keys     conversationListKeyMap
    width    int
    height   int
    selected string // Currently selected conversation ID
}

func NewConversationList(conversations []db.Conversation) ConversationListModel {
    // Convert conversations to list items
    items := make([]list.Item, len(conversations))
    for i, conv := range conversations {
        title := conv.Title
        if title == "" {
            title = fmt.Sprintf("Conversation %d", i+1)
        }
        items[i] = conversationItem{
            id:        conv.ID,
            title:     title,
            createdAt: conv.CreatedAt,
        }
    }

    // Create new list
    l := list.New(items, list.NewDefaultDelegate(), 0, 0)
    l.Title = "Conversations"
    l.SetShowHelp(true)
    l.SetFilteringEnabled(true)
    l.Styles.Title = lipgloss.NewStyle().
        Foreground(lipgloss.Color("99")).
        Bold(true).
        Padding(0, 0, 1, 2)

    return ConversationListModel{
        list: l,
        keys: conversationListKeys,
    }
}

func (m ConversationListModel) Init() tea.Cmd {
    return nil
}

func (m ConversationListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.list.SetSize(msg.Width, msg.Height)
        return m, nil

    case tea.KeyMsg:
        if msg.String() == "esc" {
            // Return to chat view
            return m, func() tea.Msg { return ChangeViewMsg(chatView) }
        }
        // Handle selection
        if key.Matches(msg, m.keys.Select) {
            item, ok := m.list.SelectedItem().(conversationItem)
            if ok {
                m.selected = item.id
                return m, func() tea.Msg { return ConversationSelectedMsg(item.id) }
            }
        }
    }

    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}

func (m ConversationListModel) View() string {
    return m.list.View()
}

// Message types for conversation list events
type ConversationSelectedMsg string 