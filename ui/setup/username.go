package setup

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UsernameInput represents the username input stage
type UsernameInput struct {
	textInput textinput.Model
	done      bool
	width     int
	height    int
}

// NewUsernameInput creates a new username input component
func NewUsernameInput() UsernameInput {
	ti := textinput.New()
	ti.Placeholder = "Enter your username"
	ti.Focus()

	return UsernameInput{
		textInput: ti,
		width:    0,
		height:   0,
	}
}

// Init initializes the username input
func (m UsernameInput) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles updates for the username input
func (m UsernameInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.textInput.Value() != "" {
				m.done = true
				return m, nil
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the username input
func (m UsernameInput) View() string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Welcome to Goatmeal!"),
			boxStyle.Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					m.textInput.View(),
				),
			),
		),
	)
}

// Done returns whether the username input is complete
func (m UsernameInput) Done() bool {
	return m.done
}

// Value returns the entered username
func (m UsernameInput) Value() string {
	return m.textInput.Value()
} 