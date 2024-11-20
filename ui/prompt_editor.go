package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PromptEditorModel struct {
	err       error
	quitting  bool
	width     int
	height    int
	prompt    string
	tempFile  string
}

type PromptEditedMsg struct {
	newPrompt string
}

func NewPromptEditor(initialPrompt string) (PromptEditorModel, error) {
	// Create temporary file
	tmpDir := filepath.Join(os.TempDir(), "goatmeal")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return PromptEditorModel{}, fmt.Errorf("could not create temp directory: %w", err)
	}

	tmpFile := filepath.Join(tmpDir, "prompt.txt")
	if err := os.WriteFile(tmpFile, []byte(initialPrompt), 0644); err != nil {
		return PromptEditorModel{}, fmt.Errorf("could not write temp file: %w", err)
	}

	return PromptEditorModel{
		prompt:   initialPrompt,
		tempFile: tmpFile,
	}, nil
}

func (m PromptEditorModel) Init() tea.Cmd {
	return m.editPrompt()
}

func (m PromptEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case error:
		m.err = msg
		return m, tea.Quit

	case PromptEditedMsg:
		// Clean up temp file
		os.Remove(m.tempFile)
		return m, func() tea.Msg { return ChangeViewMsg(systemPromptView) }
	}

	return m, nil
}

func (m PromptEditorModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	style := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width)

	return style.Render("Opening editor...\nPress ctrl+c to cancel")
}

func (m PromptEditorModel) editPrompt() tea.Cmd {
	return func() tea.Msg {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nvim"
		}

		cmd := exec.Command(editor, m.tempFile)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run the editor
		if err := cmd.Run(); err != nil {
			return err
		}

		// Read the edited content
		content, err := os.ReadFile(m.tempFile)
		if err != nil {
			return err
		}

		return PromptEditedMsg{newPrompt: string(content)}
	}
} 