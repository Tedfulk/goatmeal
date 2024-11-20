package ui

import (
	"goatmeal/config"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type inputKeyMap struct {
	Send     key.Binding
	NewLine  key.Binding
	Quit     key.Binding
	NewChat  key.Binding
}

var inputKeys = inputKeyMap{
	Send: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "send"),
	),
	NewLine: key.NewBinding(
		key.WithKeys("shift+enter"),
		key.WithHelp("shift+enter", "new line"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c/esc", "quit"),
	),
	NewChat: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "new chat"),
	),
}

type InputModel struct {
	textarea    textarea.Model
	keyMap      inputKeyMap
	width       int
	height      int
	placeholder string
}

func NewInput(colors config.ThemeColors) InputModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.ShowLineNumbers = false
	ta.SetWidth(30)
	ta.SetHeight(2)

	// Style the textarea with no bottom padding/margin and theme colors
	ta.FocusedStyle.Base = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colors.UserBubble)). // Use theme color for border
		BorderBackground(lipgloss.Color("0")).  // Black background
		Padding(0).                            // Remove all padding
		MarginTop(0).                          // Remove top margin
		MarginBottom(0).                       // Remove bottom margin
		MarginLeft(0).                         // Remove left margin
		MarginRight(0)                         // Remove right margin

	ta.BlurredStyle.Base = ta.FocusedStyle.Base

	// Style the text with theme colors
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.UserText)) // Theme text color
	ta.BlurredStyle.Text = ta.FocusedStyle.Text

	// Style the placeholder
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray placeholder
	ta.BlurredStyle.Placeholder = ta.FocusedStyle.Placeholder

	return InputModel{
		textarea:    ta,
		keyMap:      inputKeys,
		placeholder: "Send a message...",
	}
}

func (m InputModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.NewChat) {
			return m, func() tea.Msg { return NewChatMsg{} }
		}

		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Send):
			if content := m.textarea.Value(); content != "" {
				// Create command to send message
				cmds = append(cmds, func() tea.Msg {
					return SendMessageMsg{content: content}
				})
				m.textarea.Reset()
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(msg.Width - 4)
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m InputModel) View() string {
	return m.textarea.View()
}
