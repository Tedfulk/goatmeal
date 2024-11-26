package ui

import (
	"strings"

	"github.com/tedfulk/goatmeal/config"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
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
    preview   viewport.Model
    previewFocused bool
}

func NewPromptSelection(cfg *config.Config) PromptSelectionModel {
    preview := viewport.New(100, 30)
    preview.Style = lipgloss.NewStyle()

    return PromptSelectionModel{
        prompts:  cfg.SystemPrompts,
        selected: 0,
        keys:     menuKeys,
        config:   cfg,
        preview:  preview,
        previewFocused: false,
    }
}

func (m PromptSelectionModel) Init() tea.Cmd {
    return nil
}

func (m PromptSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" {
            return m, tea.Quit
        }

        if msg.String() == "esc" {
            return m, func() tea.Msg { return ChangeViewMsg(systemPromptView) }
        }

        if msg.String() == "tab" {
            m.previewFocused = !m.previewFocused
            return m, nil
        }

        if m.previewFocused {
            switch msg.String() {
            case "up", "k":
                m.preview.LineUp(1)
            case "down", "j":
                m.preview.LineDown(1)
            }
            return m, nil
        }

        switch {
        case key.Matches(msg, m.keys.Up):
            m.selected--
            if m.selected < 0 {
                m.selected = len(m.prompts) - 1
            }
            m.updatePreview()

        case key.Matches(msg, m.keys.Down):
            m.selected++
            if m.selected >= len(m.prompts) {
                m.selected = 0
            }
            m.updatePreview()

        case key.Matches(msg, m.keys.Select):
            selectedPrompt := m.prompts[m.selected]
            return m, func() tea.Msg {
                return SystemPromptSelectedMsg{prompt: selectedPrompt}
            }
        }

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        
        m.preview.Width = m.width - 44
        m.preview.Height = m.height - 8
        m.updatePreview()
    }
    return m, cmd
}

func (m *PromptSelectionModel) updatePreview() {
    if len(m.prompts) > 0 {
        m.preview.SetContent(m.prompts[m.selected])
    }
}

func (m PromptSelectionModel) View() string {
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

    listStyle := lipgloss.NewStyle().
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("99")).
        Padding(1, 2).
        Width(40)

    previewStyle := lipgloss.NewStyle().
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("99")).
        Padding(1, 2)

    if m.previewFocused {
        previewStyle = previewStyle.BorderForeground(lipgloss.Color("99"))
        listStyle = listStyle.BorderForeground(lipgloss.Color("246"))
    } else {
        previewStyle = previewStyle.BorderForeground(lipgloss.Color("246"))
        listStyle = listStyle.BorderForeground(lipgloss.Color("99"))
    }
    
    var listItems string
    for i, prompt := range m.prompts {
        words := strings.Fields(prompt)
        preview := strings.Join(words[:min(5, len(words))], " ") + "..."
        
        if i == m.selected {
            listItems += selectedStyle.Render("▸ " + preview) + "\n"
        } else {
            listItems += normalStyle.Render("  " + preview) + "\n"
        }
    }

    list := listStyle.Render(listItems)
    preview := previewStyle.Render(m.preview.View())

    panels := lipgloss.JoinHorizontal(
        lipgloss.Top,
        list,
        "  ",
        preview,
    )

    helpText := "↑/↓: navigate • tab: switch focus • enter: select • esc: back • q: quit"
    if m.previewFocused {
        helpText = "↑/↓: scroll • up/down: page scroll • tab: switch focus • esc: back • q: quit"
    }

    content := lipgloss.JoinVertical(
        lipgloss.Center,
        titleStyle.Render("Select System Prompt"),
        "",
        panels,
        "",
        normalStyle.Render(helpText),
    )

    return lipgloss.Place(
        m.width,
        m.height,
        lipgloss.Center,
        lipgloss.Center,
        content,
    )
}
