package ui

import (
	"goatmeal/config"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PromptSelectionModel struct {
    prompts   []string
    selected  int
    keys      menuKeyMap
    width     int
    height    int
    config    *config.Config
}

func NewPromptSelection(cfg *config.Config) PromptSelectionModel {
    return PromptSelectionModel{
        prompts:  cfg.SystemPrompts,
        selected: 0,
        keys:     menuKeys,
        config:   cfg,
    }
}

func (m PromptSelectionModel) Init() tea.Cmd {
    return nil
}

func (m PromptSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" {
            return m, tea.Quit
        }

        if msg.String() == "esc" {
            return m, func() tea.Msg { return ChangeViewMsg(systemPromptView) }
        }

        switch {
        case key.Matches(msg, m.keys.Up):
            m.selected--
            if m.selected < 0 {
                m.selected = len(m.prompts) - 1
            }

        case key.Matches(msg, m.keys.Down):
            m.selected++
            if m.selected >= len(m.prompts) {
                m.selected = 0
            }

        case key.Matches(msg, m.keys.Select):
            selectedPrompt := m.prompts[m.selected]
            return m, func() tea.Msg {
                return SystemPromptSelectedMsg{prompt: selectedPrompt}
            }
        }
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    return m, nil
}

func (m PromptSelectionModel) View() string {
    // Create styles
    titleStyle := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("99")).
        Padding(1, 0).
        Align(lipgloss.Center)

    selectedStyle := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("99"))

    normalStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("246"))

    // Create list panel (left side)
    listStyle := lipgloss.NewStyle().
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("99")).
        Padding(1, 2).
        Width(40)  // Fixed width for list panel

    // Create preview panel (right side)
    previewStyle := lipgloss.NewStyle().
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("99")).
        Padding(1, 2).
        Width(m.width - 44)  // Remaining width minus list panel and spacing

    // Build list items
    var listItems string
    for i, prompt := range m.prompts {
        // Get first 5 words for preview
        words := strings.Fields(prompt)
        preview := strings.Join(words[:min(5, len(words))], " ") + "..."
        
        if i == m.selected {
            listItems += selectedStyle.Render("▸ " + preview) + "\n"
        } else {
            listItems += normalStyle.Render("  " + preview) + "\n"
        }
    }

    // Create preview content
    var previewContent string
    if len(m.prompts) > 0 {
        previewContent = m.prompts[m.selected]
    }

    // Create panels
    list := listStyle.Render(listItems)
    preview := previewStyle.Render(previewContent)

    // Join panels horizontally
    panels := lipgloss.JoinHorizontal(
        lipgloss.Top,
        list,
        "  ",  // 2-space gap between panels
        preview,
    )

    // Create the complete view
    content := lipgloss.JoinVertical(
        lipgloss.Center,
        titleStyle.Render("Select System Prompt"),
        "",
        panels,
        "",
        normalStyle.Render("↑/↓: navigate • enter: select • esc: back • q: quit"),
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
