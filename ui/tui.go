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

// Add view states
type viewState int

const (
	chatView viewState = iota
	menuView
)

// MainModel combines all the components of our chat UI
type MainModel struct {
	chat        ChatModel
	input       InputModel
	menu        MenuModel
	currentView viewState
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
		chat:        chat,
		input:      NewInput(),
		menu:       NewMenu(),
		currentView: chatView,
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
		
		if m.currentView == menuView {
			var menuModel tea.Model
			menuModel, cmd = m.menu.Update(msg)
			m.menu = menuModel.(MenuModel)
			return m, cmd
		}
		
		// Update chat components with new size
		m.chat, cmd = m.chat.Update(msg)
		cmds = append(cmds, cmd)
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		// Global key handlers
		switch msg.String() {
		case "shift+tab":
			// Toggle between chat and menu views
			if m.currentView == chatView {
				m.currentView = menuView
				m.input.textarea.Blur()
			} else {
				m.currentView = chatView
				m.input.textarea.Focus()
			}
			return m, nil
		}

		// Handle view-specific updates
		if m.currentView == menuView {
			var menuModel tea.Model
			menuModel, cmd = m.menu.Update(msg)
			m.menu = menuModel.(MenuModel)
			
			// Check if menu is quitting
			if m.menu.quitting {
				m.currentView = chatView
				m.input.textarea.Focus()
			}
			return m, cmd
		}

		// Chat view key handlers
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
				
				m.chat.AddMessage(userMsg)
				m.loading = true
				
				return m, tea.Batch(
					m.getAIResponse(userMsg),
					m.spinner.Tick,
				)
			}
		}

		// Handle tab for focus switching in chat view
		if msg.String() == "tab" {
			if m.input.textarea.Focused() {
				m.input.textarea.Blur()
			} else {
				m.input.textarea.Focus()
			}
			return m, nil
		}

		// Pass other messages to appropriate component
		if !m.input.textarea.Focused() {
			m.chat, cmd = m.chat.Update(msg)
		} else {
			m.input, cmd = m.input.Update(msg)
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	if m.currentView == menuView {
		// Center the entire menu view in the terminal
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			m.menu.View(),
		)
	}

	// Show chat view
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.chat.View(),
		m.input.View(),
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