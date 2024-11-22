package ui

import (
	"fmt"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/db"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ConversationListModel struct {
	table         table.Model
	preview       viewport.Model
	selectedConv  string
	width         int
	height        int
	colors        config.ThemeColors
	db            db.ChatDB
	renderer      *glamour.TermRenderer
	conversations []db.Conversation
	previewFocused bool
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62"))

func NewConversationList(conversations []db.Conversation, colors config.ThemeColors, database db.ChatDB) ConversationListModel {
	// Create table rows from conversations
	rows := make([]table.Row, len(conversations))
	for i, conv := range conversations {
		title := conv.Title
		if title == "" {
			title = fmt.Sprintf("Conversation %d", i+1)
		}
		rows[i] = table.Row{title, conv.ID}
	}

	// Define table columns
	columns := []table.Column{
		{Title: "Conversations", Width: 30},
		{Title: "ID", Width: 0}, // Hidden column for ID
	}

	// Initialize table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(34),
	)

	// Style the table
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(colors.MenuTitle)).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(colors.MenuSelected)).
		Bold(true)
	t.SetStyles(s)

	// Initialize viewport for preview
	vp := viewport.New(102, 35)
	vp.Style = lipgloss.NewStyle()

	// Initialize glamour renderer
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(90),
	)

	return ConversationListModel{
		table:          t,
		preview:        vp,
		colors:         colors,
		db:             database,
		renderer:       renderer,
		conversations:  conversations,
		previewFocused: false,
	}
}

func (m ConversationListModel) Init() tea.Cmd {
	return nil
}

func (m ConversationListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Update table size to fit the left panel
		m.table.SetWidth(m.width/3)
		m.table.SetHeight(m.height - 4)
		
		// Update preview size
		m.preview.Width = (m.width * 2 / 3) - 4
		m.preview.Height = m.height - 4

		// Force preview update after resize
		m.updatePreview()
		
		return m, nil

	case tea.KeyMsg:
		// Handle quitting the program
		if msg.String() == "q" {
			return m, tea.Quit
		}

		// Handle tab key to switch focus
		if msg.String() == "tab" {
			m.previewFocused = !m.previewFocused
			return m, nil
		}

		// Handle other keys based on focus
		if m.previewFocused {
			// When preview is focused, send key events to viewport
			var cmd tea.Cmd
			m.preview, cmd = m.preview.Update(msg)
			return m, cmd
		} else {
			// When table is focused, handle table navigation
			switch msg.String() {
			case "esc":
				return m, func() tea.Msg { return ChangeViewMsg(chatView) }
			case "up", "down":
				// Update table first
				m.table, cmd = m.table.Update(msg)
				// Then force preview update
				m.updatePreview()
				return m, cmd
			default:
				m.table, cmd = m.table.Update(msg)
				m.updatePreview()
				return m, cmd
			}
		}
	}

	return m, cmd
}

func (m *ConversationListModel) updatePreview() {
	selectedRow := m.table.SelectedRow()
	if len(selectedRow) < 2 {
		return
	}

	convID := selectedRow[1] // Get ID from second column

	// First verify the conversation exists
	conv, err := m.db.GetConversation(convID)
	if err != nil {
		m.preview.SetContent("Error loading conversation")
		return
	}
	if conv == nil {
		m.preview.SetContent("Conversation not found")
		return
	}

	messages, err := m.db.GetMessages(convID)
	if err != nil {
		m.preview.SetContent("Error loading messages")
		return
	}

	if len(messages) == 0 {
		m.preview.SetContent("No messages in this conversation")
		return
	}

	// Initialize the title caser
	caser := cases.Title(language.English)

	// Convert messages to markdown format
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", conv.Title))
	
	for _, msg := range messages {
		sb.WriteString(fmt.Sprintf("## %s\n", caser.String(msg.Role)))
		sb.WriteString(fmt.Sprintf("%s\n\n", msg.Content))
	}

	// Render markdown
	content, err := m.renderer.Render(sb.String())
	if err != nil {
		m.preview.SetContent("Error rendering messages")
		return
	}

	m.preview.SetContent(content)
	m.preview.GotoTop()
}

func (m ConversationListModel) View() string {
	// Create a container for the table that stays in the top left
	tableStyle := lipgloss.NewStyle().
		Width(m.width/3).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.colors.MenuBorder))

	// Create a container for the preview that fills the right side
	previewStyle := lipgloss.NewStyle().
		Width((m.width * 2 / 3) - 4).
		Height(m.height - 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.colors.MenuBorder))

	// Apply focus styles
	if m.previewFocused {
		previewStyle = previewStyle.BorderForeground(lipgloss.Color(m.colors.MenuSelected))
	} else {
		tableStyle = tableStyle.BorderForeground(lipgloss.Color(m.colors.MenuSelected))
	}

	tableContainer := tableStyle.Render(m.table.View())
	previewContainer := previewStyle.Render(m.preview.View())

	// Join the panels horizontally with a small gap
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		tableContainer,
		lipgloss.NewStyle().Width(2).Render(""), // Add gap between panels
		previewContainer,
	)
}

// Message types for conversation list events
type ConversationSelectedMsg string

// Add this method to initialize the preview with the first conversation
func (m *ConversationListModel) initializePreview() {
	// If there are conversations, show the first one
	if len(m.table.Rows()) > 0 {
		// Get the first row's conversation ID
		firstRow := m.table.Rows()[0]
		if len(firstRow) > 1 {
			convID := firstRow[1]
			
			// Set the table selection to the first row
			m.table.SetCursor(0)
			
			// Update the preview with the first conversation
			m.selectedConv = convID
			m.updatePreview()
			
			// Make sure the viewport is at the top
			m.preview.GotoTop()
		}
	}
}