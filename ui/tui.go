package ui

import (
	"fmt"
	"goatmeal/api"
	"goatmeal/config"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Message types for our Update function
type errMsg error
type aiResponseMsg api.Message

// MainModel combines all the components of our chat UI
type MainModel struct {
	chat        ChatModel
	input       InputModel
	spinner     spinner.Model
	loading     bool
	height      int
	width       int
	err         error
	quitting    bool
	groqClient  *api.GroqClient
	config      *config.Config
}

func NewMainModel() (MainModel, error) {
	// Load the configuration
	config, err := config.LoadConfig()
	if err != nil {
		return MainModel{}, fmt.Errorf("error loading config: %w", err)
	}

	// Initialize the Groq client with the loaded config
	groqClient, err := api.NewGroqClient(config)
	if err != nil {
		return MainModel{}, fmt.Errorf("error creating Groq client: %w", err)
	}

	// Initialize the spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Initialize chat with config for proper theming
	chat, err := NewChat(config)
	if err != nil {
		return MainModel{}, fmt.Errorf("error creating chat view: %w", err)
	}

	return MainModel{
		chat:       chat,
		input:      NewInput(),
		spinner:    s,
		groqClient: groqClient,
		config:     config,
	}, nil
}

func (m MainModel) Init() tea.Cmd {
	return tea.Batch(
		m.input.Init(),
		m.chat.Init(),
		m.spinner.Tick,
	)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		
		// Update both components with new size
		m.chat, cmd = m.chat.Update(msg)
		cmds = append(cmds, cmd)
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)

	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case aiResponseMsg:
		m.loading = false
		m.chat.AddMessage(api.Message(msg))
		m.input.textarea.Reset()
		return m, nil

	case errMsg:
		m.loading = false
		m.err = msg
		errorMessage := api.Message{
			Role:      "system",
			Content:   fmt.Sprintf("Error: %v", msg),
			Timestamp: time.Now(),
		}
		m.chat.AddMessage(errorMessage)
		return m, nil

	case tea.KeyMsg:
		// Handle tab key for both scroll and input modes
		if msg.String() == "tab" {
			if m.input.textarea.Focused() {
				m.input.textarea.Blur()
			} else {
				m.input.textarea.Focus()
			}
			return m, nil
		}

		// Handle scrolling when not focused on input
		if !m.input.textarea.Focused() {
			var chatCmd tea.Cmd
			m.chat, chatCmd = m.chat.Update(msg)
			return m, chatCmd
		}

		// Handle global keypresses
		switch {
		case key.Matches(msg, m.input.keyMap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.input.keyMap.Send):
			if m.input.textarea.Value() != "" {
				userMsg := api.Message{
					Role:      "user",
					Content:   m.input.textarea.Value(),
					Timestamp: time.Now(),
				}
				
				// Add user message to chat
				m.chat.AddMessage(userMsg)
				m.loading = true
				
				// Create a command to get the AI response
				return m, tea.Batch(
					m.getAIResponse(userMsg),
					m.spinner.Tick,
				)
			}
		}

		// If input is focused, pass events to input
		if m.input.textarea.Focused() {
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	// Combine chat and input views without the top border
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.chat.View(),
		m.input.View(), // Remove the border styling wrapper
	)
}

// Helper function to create a command for getting AI response
func (m MainModel) getAIResponse(userMsg api.Message) tea.Cmd {
	return func() tea.Msg {
		resp, err := m.groqClient.SendMessage(userMsg.Content, m.chat.GetMessages())
		if err != nil {
			return errMsg(err)
		}
		return aiResponseMsg(*resp)
	}
}