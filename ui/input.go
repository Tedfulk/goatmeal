package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tedfulk/goatmeal/ui/theme"
)

// Input represents the user input component
type Input struct {
	textInput textinput.Model
	Width     int
}

// NewInput creates a new input component
func NewInput() Input {
	ti := textinput.New()
	ti.Placeholder = "Type your message..."
	ti.Focus()
	ti.Width = 100
	ti.Prompt = " > "

	return Input{
		textInput: ti,
	}
}

// Update handles input updates
func (i Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	i.textInput, cmd = i.textInput.Update(msg)
	return i, cmd
}

// View renders the input component
func (i Input) View() string {
	return theme.BaseStyle.Input.
		Width(i.Width).
		BorderForeground(theme.CurrentTheme.Primary.GetColor()).
		Render(i.textInput.View())
}

// Value returns the current input value
func (i Input) Value() string {
	return i.textInput.Value()
}

// Reset clears the input
func (i *Input) Reset() {
	i.textInput.Reset()
} 