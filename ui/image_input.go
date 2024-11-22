package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
)

type ImageInputModel struct {
	imageInput   textinput.Model
	contextInput textinput.Model
	width        int
	height       int
	err          error
	focused      bool
	colors       config.ThemeColors
}

type ImageInputMsg struct {
	imagePath string
	context   string
}

func NewImageInput(colors config.ThemeColors) ImageInputModel {
	contextInput := textinput.New()
	contextInput.Placeholder = "Add additional context"
	contextInput.Focus()
	contextInput.CharLimit = 256
	contextInput.Width = 70

	imageInput := textinput.New()
	imageInput.Placeholder = "Enter Image URL or Path"
	imageInput.CharLimit = 256
	imageInput.Width = 70

	return ImageInputModel{
		imageInput:   imageInput,
		contextInput: contextInput,
		focused:      true,
		colors:       colors,
	}
}

func (m ImageInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ImageInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return ChangeViewMsg(chatView) }
		case "enter":
			if m.imageInput.Value() != "" && m.contextInput.Value() != "" {
				return m, func() tea.Msg {
					return ImageInputMsg{
						imagePath: m.imageInput.Value(),
						context:   m.contextInput.Value(),
					}
				}
			}
		case "tab":
			m.focused = !m.focused
			if m.focused {
				m.contextInput.Focus()
				m.imageInput.Blur()
			} else {
				m.imageInput.Focus()
				m.contextInput.Blur()
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	if m.focused {
		m.contextInput, cmd = m.contextInput.Update(msg)
	} else {
		m.imageInput, cmd = m.imageInput.Update(msg)
	}
	return m, cmd
}

func (m ImageInputModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colors.UserText)).
		Padding(0, 1).
		MarginBottom(0).
		Width(80).
		Align(lipgloss.Center)

	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.colors.UserBubble)).
		Padding(1).
		Width(80)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colors.MenuNormal)).
		Align(lipgloss.Center)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("Image Input"),
		baseStyle.Render(m.contextInput.View()),
		baseStyle.Render(m.imageInput.View()),
		helpStyle.Render("Press Tab to switch focus • Enter to submit • Esc to cancel"),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
} 