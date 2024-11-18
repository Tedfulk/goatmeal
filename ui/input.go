package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Send     key.Binding
	NewLine  key.Binding
	Quit     key.Binding
}

var keys = keyMap{
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
}

type InputModel struct {
	textarea    textarea.Model
	keyMap      keyMap
	width       int
	height      int
	placeholder string
}

func NewInput() InputModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.ShowLineNumbers = false
	ta.SetWidth(30)
	ta.SetHeight(4)

	// Style the textarea
	ta.FocusedStyle.Base = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("7")). // White border
		BorderBackground(lipgloss.Color("0")).  // Black background
		MarginBottom(0)                        // Remove bottom margin

	ta.BlurredStyle.Base = ta.FocusedStyle.Base

	// Style the text
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("7")) // White text
	ta.BlurredStyle.Text = ta.FocusedStyle.Text

	// Style the placeholder
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray placeholder
	ta.BlurredStyle.Placeholder = ta.FocusedStyle.Placeholder

	return InputModel{
		textarea:    ta,
		keyMap:      keys,
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
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Send):
			if m.textarea.Value() != "" {
				// Handle send message here
				m.textarea.Reset()
				return m, nil
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
