package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type APIKeyEditorModel struct {
	textInput textinput.Model
	width     int
	height    int
	err       error
}

type APIKeyUpdatedMsg struct {
	newKey string
}

func NewAPIKeyEditor(currentKey string) APIKeyEditorModel {
	ti := textinput.New()
	ti.Placeholder = "Enter API Key"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 70

	return APIKeyEditorModel{
		textInput: ti,
	}
}

func (m APIKeyEditorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m APIKeyEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return ChangeViewMsg(settingsView) }
		case "enter":
			if m.textInput.Value() != "" {
				return m, func() tea.Msg {
					return APIKeyUpdatedMsg{newKey: m.textInput.Value()}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m APIKeyEditorModel) View() string {
	// Create API Key title style
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("white")).
		Padding(0, 1).
		MarginBottom(0).
		Width(80).
		Align(lipgloss.Center)

	// Create the box with input
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("white")).
		Padding(1).
		Width(80)

	// Help text style
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center)

	// Combine all elements vertically
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("API Key"),
		baseStyle.Render(m.textInput.View()),
		helpStyle.Render("Press Enter to submit • Esc to cancel"),
	)

	// Center everything in the terminal
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
} 