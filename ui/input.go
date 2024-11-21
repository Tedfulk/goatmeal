package ui

import (
	"github.com/tedfulk/goatmeal/config"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type inputKeyMap struct {
	Send         key.Binding
	NewLine      key.Binding
	Quit         key.Binding
	NewChat      key.Binding
	ThemeSelector key.Binding
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
		key.WithKeys("ctrl+t"),
		key.WithHelp("ctrl+t", "new chat"),
	),
	ThemeSelector: key.NewBinding(
		key.WithKeys(";"),
		key.WithHelp(";", "go to theme selector"),
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
	// Create a clean textarea first
	ta := textarea.New()
	
	// Create a completely clean base style by unsetting everything
	baseStyle := lipgloss.NewStyle().
		UnsetAlign().
		UnsetBackground().
		UnsetBold().
		UnsetBorderStyle().
		UnsetForeground().
		UnsetHeight().
		UnsetWidth().
		UnsetPadding().
		UnsetMargins()

	// Now build our styles on the clean base
	inputStyle := baseStyle.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colors.UserBubble)).
		Background(lipgloss.Color("0")).
		Padding(0).
		Margin(0)

	// Create clean text style
	textStyle := lipgloss.NewStyle().
		UnsetBorderStyle().
		Foreground(lipgloss.Color(colors.UserText)).
		Background(lipgloss.Color("0"))

	// Create clean placeholder style
	placeholderStyle := lipgloss.NewStyle().
		UnsetBorderStyle().
		Foreground(lipgloss.Color("240")).
		Background(lipgloss.Color("0"))

	// Create clean cursor style
	cursorLineStyle := lipgloss.NewStyle().
		UnsetBorderStyle().
		Background(lipgloss.Color("0"))

	// Set basic properties
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.SetWidth(138)
	ta.SetHeight(2)
	
	// Apply our clean styles
	ta.BlurredStyle.Base = inputStyle
	ta.BlurredStyle.Text = textStyle
	ta.BlurredStyle.Placeholder = placeholderStyle
	
	ta.FocusedStyle.Base = inputStyle
	ta.FocusedStyle.Text = textStyle
	ta.FocusedStyle.Placeholder = placeholderStyle
	ta.FocusedStyle.CursorLine = cursorLineStyle

	// Set placeholder after styles
	ta.Placeholder = "Send a message..."
	
	// Reset and focus
	ta.Reset()
	ta.Focus()

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

func (m *InputModel) Reset() {
	// Create a fresh textarea
	ta := textarea.New()
	ta.Placeholder = m.placeholder
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.Focus()

	// Copy dimensions
	ta.SetWidth(m.textarea.Width())
	ta.SetHeight(m.textarea.Height())

	// Copy styles
	ta.FocusedStyle = m.textarea.FocusedStyle
	ta.BlurredStyle = m.textarea.BlurredStyle

	// Ensure cursor line has proper background
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().
		Background(lipgloss.Color("0"))

	// Explicitly clear content
	ta.SetValue("")

	// Reset internal state
	ta.Reset()

	// Replace the old textarea
	m.textarea = ta

	// Force focus to ensure proper state
	m.textarea.Focus()
}
