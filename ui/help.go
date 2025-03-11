package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/ui/theme"
)

type HelpView struct {
	viewport viewport.Model
	width    int
	height   int
}

func NewHelpView() *HelpView {
	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.CurrentTheme.Primary.GetColor()).
		Padding(1, 2)

	// Create the help content
	content := `# Keyboard Shortcuts

## Navigation
* **ctrl+t**: Start a new conversation
* **ctrl+l**: View conversation list
* **ctrl+s**: Open settings menu
* **ctrl+c**: Quit application
* **esc**: Go back/close current view

## Chat Interface
* **?**: Toggle menu
* **enter**: Send message
* **/o[n]**: Open message number 'n' in editor (e.g., /o1)
* **/c[n]**: Copy message number 'n' to clipboard (e.g., /c1)
* **/b[n]**: Copy code block number 'n' to clipboard (e.g., /b1)
* **/s[n]**: Speak message number 'n' (e.g., /s1)
* **ctrl+q**: Stop current speech playback

## Web Search Commands
* **/web query**: Search for information
* **/web query +domain.com**: Search with specific domain
* **/webe query**: Enhanced web search with AI optimization
* **/webe query +domain.com**: Enhanced domain-specific search
* **/epq query**: Enhanced programming query with AI optimization

Enhanced search mode (🔍+) uses AI to optimize your search queries for better results. It adds context, clarity, and relevant details while maintaining a concise format.

Enhanced programming query mode (💻+) uses AI to optimize programming-related queries by adding specificity about languages, frameworks, and technical requirements.

## Conversation List
* **tab**: Switch focus between list and messages
* **ctrl+d**: Delete selected conversation
* **ctrl+e**: Export conversation as JSON
* **esc**: Return to chat

## Settings Menu
* **enter**: Select option
* **esc**: Return to previous menu

## General
* **ctrl+c, q**: Quit application
`

	// Render the markdown content
	rendered, err := glamour.Render(content, "dark")
	if err != nil {
		rendered = content // Fallback to plain text if rendering fails
	}
	vp.SetContent(rendered)

	return &HelpView{
		viewport: vp,
	}
}

func (h *HelpView) Update(msg tea.Msg) (*HelpView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return h, func() tea.Msg {
				return SetViewMsg{view: "chat"}
			}
		}
	}

	h.viewport, cmd = h.viewport.Update(msg)
	return h, cmd
}

func (h *HelpView) View() string {
	return h.viewport.View()
}

func (h *HelpView) SetSize(width, height int) {
	h.width = width
	h.height = height
	h.viewport.Width = width - 4
	h.viewport.Height = height - 4
}
