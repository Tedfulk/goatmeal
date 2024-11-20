package ui

import (
	"fmt"
	"goatmeal/api"
	"goatmeal/config"
	"goatmeal/db"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Message types for our Update function
type errMsg error
type aiResponseMsg api.Message
type newConversationMsg struct{}
type NewChatMsg struct{}
type titleResponseMsg string
type SendMessageMsg struct {
	content string
}

// Add view states
type viewState int

const (
	chatView viewState = iota
	menuView
	conversationListView
	helpView
	settingsView
	systemPromptView
	promptEditorView
	apiKeyEditorView
	themeSelectorView
	promptSelectionView
)

// MainModel combines all the components of our chat UI
type MainModel struct {
	chat                ChatModel
	input               InputModel
	menu                MenuModel
	currentView         viewState
	spinner             spinner.Model
	loading             bool
	height              int
	width               int
	err                 error
	quitting            bool
	groqClient          *api.GroqClient
	config              *config.Config
	currentConversation string        // Track current conversation ID
	conversations       []db.Conversation // Change this line to use db.Conversation
	db                  db.ChatDB      // Database interface
	conversationList    ConversationListModel
	lastView            viewState  // Track the last view for returning from conversation list
	help                HelpModel
	settings            SettingsModel
	systemPrompt        SystemPromptModel
	promptEditor        PromptEditorModel
	apiKeyEditor        APIKeyEditorModel
	themeSelector       ThemeSelectorModel
	promptSelection     PromptSelectionModel
	conversationCreated  bool // Add a flag to track if a conversation is created
	previewFocused      bool // Track if the preview is focused
}

func NewMainModel(db db.ChatDB) (MainModel, error) {
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

	// Initialize chat without a conversation ID
	chat, err := NewChat(config, db, "")
	if err != nil {
		return MainModel{}, fmt.Errorf("error creating chat view: %w", err)
	}

	// Initialize the spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Get conversations for the list
	conversations, err := db.ListConversations()
	if err != nil {
		return MainModel{}, fmt.Errorf("error listing conversations: %w", err)
	}

	colors := config.GetThemeColors()

	return MainModel{
		chat:                chat,
		input:               NewInput(colors),
		menu:                NewMenu(colors),
		currentView:         chatView,
		spinner:             s,
		groqClient:          groqClient,
		config:              config,
		db:                  db,
		conversationList:    NewConversationList(conversations, colors, db),
		lastView:            chatView,
		help:                NewHelp(colors),
		settings:            NewSettings(colors),
		systemPrompt:        NewSystemPromptMenu(config),
		apiKeyEditor:        NewAPIKeyEditor(config.APIKey),
		themeSelector:       NewThemeSelector(colors),
		promptSelection:     NewPromptSelection(config),
		conversationCreated: false, // Initialize the flag as false
		previewFocused:      false,  // Initialize preview focus state
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
		// Handle ESC in chat view to quit
		if msg.String() == "esc" && m.currentView == chatView {
			m.quitting = true
			return m, tea.Quit
		}

		// Global key handlers
		switch msg.String() {
		case "shift+tab":
			if m.currentView == chatView {
				m.currentView = menuView
				m.input.textarea.Blur()
			} else {
				m.currentView = chatView
				m.input.textarea.Focus()
			}
			return m, nil
		case "tab":
			if m.currentView == chatView {
				if m.input.textarea.Focused() {
					m.input.textarea.Blur()
					m.chat.focused = true
				} else {
					m.input.textarea.Focus()
					m.chat.focused = false
				}
				return m, nil
			}
		case "ctrl+l":
			m.lastView = m.currentView
			m.currentView = conversationListView
			m.conversationList.initializePreview()
			conversations, err := m.db.ListConversations()
			if err == nil {
				m.conversations = conversations
				m.conversationList = NewConversationList(conversations, m.config.GetThemeColors(), m.db)
				m.conversationList.initializePreview()
			}
			return m, nil
		}
		// Handle view-specific updates
		switch m.currentView {
		case apiKeyEditorView:
			var editorModel tea.Model
			editorModel, cmd = m.apiKeyEditor.Update(msg)
			m.apiKeyEditor = editorModel.(APIKeyEditorModel)
			return m, cmd
		case promptEditorView:
			var editorModel tea.Model
			editorModel, cmd = m.promptEditor.Update(msg)
			m.promptEditor = editorModel.(PromptEditorModel)
			return m, cmd
		case conversationListView:
			var listModel tea.Model
			listModel, cmd = m.conversationList.Update(msg)
			m.conversationList = listModel.(ConversationListModel)
			return m, cmd
		case menuView:
			var menuModel tea.Model
			menuModel, cmd = m.menu.Update(msg)
			m.menu = menuModel.(MenuModel)
			return m, cmd
		case settingsView:
			var settingsModel tea.Model
			settingsModel, cmd = m.settings.Update(msg)
			m.settings = settingsModel.(SettingsModel)
			return m, cmd
		case systemPromptView:
			var promptModel tea.Model
			promptModel, cmd = m.systemPrompt.Update(msg)
			m.systemPrompt = promptModel.(SystemPromptModel)
			return m, cmd
		case helpView:
			if msg.String() == "esc" {
				m.currentView = menuView
				return m, nil
			}
			var helpModel tea.Model
			helpModel, cmd = m.help.Update(msg)
			m.help = helpModel.(HelpModel)
			return m, cmd
		case themeSelectorView:
			var selectorModel tea.Model
			selectorModel, cmd = m.themeSelector.Update(msg)
			m.themeSelector = selectorModel.(ThemeSelectorModel)
			return m, cmd
		case promptSelectionView:
			var selectionModel tea.Model
			selectionModel, cmd = m.promptSelection.Update(msg)
			m.promptSelection = selectionModel.(PromptSelectionModel)
			return m, cmd
		default:
			// Chat view updates
			if m.input.textarea.Focused() {
				// Only update input when focused
				m.input, cmd = m.input.Update(msg)
				return m, cmd
			} else {
				// Handle capslock in chat view when input is not focused
				if msg.String() == "capslock" {
					newChat, err := NewChat(m.config, m.db, "")
					if err != nil {
						m.err = err
						return m, nil
					}

					// Update all relevant state
					m.currentConversation = ""
					m.chat = newChat
					m.input.textarea.Reset()
					m.conversationCreated = false
					m.currentView = chatView
					m.input.textarea.Focus()

					return m, nil
				}
				// Update chat view when input is not focused
				m.chat, cmd = m.chat.Update(msg)
				return m, cmd
			}
		}

	case aiResponseMsg:
		// Add the AI's response to the chat
		m.chat.AddMessage(api.Message(msg))
		m.loading = false
		m.input.textarea.Reset()
		return m, nil

	case errMsg:
		// Handle any errors from the API
		m.err = msg
		m.loading = false
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case NewChatMsg:
		// Reset the chat view for a new conversation without creating a database record
		newChat, err := NewChat(m.config, m.db, "")
		if err != nil {
			m.err = err
			return m, nil
		}

		// Update all relevant state
		m.currentConversation = "" // No conversation ID yet
		m.chat = newChat
		m.input.textarea.Reset()
		m.conversationCreated = false // Important: mark that no conversation exists in DB yet
		m.currentView = chatView     // Switch back to chat view
		m.input.textarea.Focus()     // Focus the input for immediate typing

		return m, nil

	case titleResponseMsg:
		// Update the conversation title in the database
		err := m.db.UpdateConversationTitle(m.currentConversation, string(msg))
		if err != nil {
			m.err = err
		}
		return m, nil

	case ConversationSelectedMsg:
		// Load the selected conversation
		convID := string(msg)
		newChat, err := NewChat(m.config, m.db, convID)
		if err != nil {
			m.err = err
			return m, nil
		}
		m.chat = newChat
		m.currentConversation = convID
		m.currentView = chatView
		return m, nil

	case ChangeViewMsg:
		m.lastView = m.currentView
		m.currentView = viewState(msg)
		return m, nil

	case SettingsMsg:
		switch msg.action {
		case EditAPIKey:
			m.currentView = apiKeyEditorView
			return m, nil
		case EditTheme:
			m.currentView = themeSelectorView
			return m, nil
		case EditSystemPrompt:
			m.lastView = m.currentView
			m.currentView = systemPromptView
			return m, nil
		}

	case SystemPromptMsg:
		switch msg.action {
		case AddPrompt:
			editor, err := NewPromptEditor("")
			if err != nil {
				m.err = err
				return m, nil
			}
			m.promptEditor = editor
			m.lastView = m.currentView
			m.currentView = promptEditorView
			return m, m.promptEditor.Init()

		case EditPrompt:
			editor, err := NewPromptEditor(msg.prompt)
			if err != nil {
				m.err = err
				return m, nil
			}
			m.promptEditor = editor
			m.lastView = m.currentView
			m.currentView = promptEditorView
			return m, m.promptEditor.Init()
		}

	case PromptEditedMsg:
		// Save the new prompt to config
		if err := m.config.SaveSystemPrompt(msg.newPrompt); err != nil {
			m.err = err
			return m, nil
		}
		// Reload the config
		newConfig, err := config.LoadConfig()
		if err != nil {
			m.err = err
			return m, nil
		}
		m.config = newConfig
		// Update the system prompt menu
		m.systemPrompt = NewSystemPromptMenu(m.config)
		return m, nil

	case APIKeyUpdatedMsg:
		// Save the new API key to config
		if err := m.config.SaveAPIKey(msg.newKey); err != nil {
			m.err = err
			return m, nil
		}
		// Reload config and update client
		newConfig, err := config.LoadConfig()
		if err != nil {
			m.err = err
			return m, nil
		}
		m.config = newConfig
		m.groqClient, err = api.NewGroqClient(newConfig)
		if err != nil {
			m.err = err
			return m, nil
		}
		// Return to settings view
		m.currentView = settingsView
		return m, nil

	case ThemeSelectedMsg:
		// Save the new theme
		if err := m.config.SaveTheme(msg.theme); err != nil {
			m.err = err
			return m, nil
		}
		
		// Reload config to apply new theme
		newConfig, err := config.LoadConfig()
		if err != nil {
			m.err = err
			return m, nil
		}
		m.config = newConfig
		
		// Get new theme colors
		colors := newConfig.GetThemeColors()
		
		// Create new chat with updated theme
		newChat, err := NewChat(newConfig, m.db, m.currentConversation)
		if err != nil {
			m.err = err
			return m, nil
		}
		
		// Copy over existing messages to the new chat model
		for _, msg := range m.chat.GetMessages() {
			if err := newChat.AddMessage(msg); err != nil {
				m.err = err
				return m, nil
			}
		}
		
		// Update all components with new colors
		m.chat = newChat
		m.input = NewInput(colors)
		m.menu = NewMenu(colors)
		m.settings = NewSettings(colors)
		m.help = NewHelp(colors)
		m.conversationList = NewConversationList(m.conversations, colors, m.db)
		m.themeSelector = NewThemeSelector(colors)
		
		// Force refresh of all components
		sizeMsg := tea.WindowSizeMsg{
			Width:  m.width,
			Height: m.height,
		}
		
		return m, func() tea.Msg { return sizeMsg }

	case SystemPromptSelectedMsg:
		if err := m.config.SetActiveSystemPrompt(msg.prompt); err != nil {
			m.err = err
			return m, nil
		}
		m.currentView = systemPromptView
		return m, nil

	case SendMessageMsg:
		// Check if a conversation has been created
		if !m.conversationCreated {
			// Create a new conversation in the database
			convID, err := m.db.CreateConversation()
			if err != nil {
				m.err = err
				return m, nil
			}

			// Update the current conversation ID
			m.currentConversation = convID
			m.chat.currentID = convID  // Make sure chat model has the ID
			m.conversationCreated = true
		}

		// Create user message
		userMsg := api.Message{
			Role:      "user",
			Content:   msg.content,
			Timestamp: time.Now(),
		}

		// Add user message to chat
		if err := m.chat.AddMessage(userMsg); err != nil {
			m.err = err
			return m, nil
		}

		// Set loading state
		m.loading = true

		// Generate title for new conversations
		var cmds []tea.Cmd
		if len(m.chat.GetMessages()) == 1 {
			cmds = append(cmds, m.generateTitle(msg.content))
		}

		// Get AI response
		cmds = append(cmds, m.getAIResponse(userMsg))
		return m, tea.Batch(cmds...)
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) centerView(content string) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m MainModel) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	switch m.currentView {
	case helpView:
		return m.centerView(m.help.View())
	case settingsView:
		return m.centerView(m.settings.View())
	case systemPromptView:
		return m.centerView(m.systemPrompt.View())
	case conversationListView:
		return m.conversationList.View()
	case menuView:
		return m.centerView(m.menu.View())
	case promptEditorView:
		return m.promptEditor.View()
	case apiKeyEditorView:
		return m.centerView(m.apiKeyEditor.View())
	case themeSelectorView:
		return m.centerView(m.themeSelector.View())
	case promptSelectionView:
		return m.promptSelection.View()
	default:
		// Remove any spacing in the vertical join
		return lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().
				MarginBottom(0).
				Render(m.chat.View()),
			lipgloss.NewStyle().
				MarginTop(0).
				MarginBottom(0).
				Render(m.input.View()),
		)
	}
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

// Add new command function
func newConversationCmd() tea.Msg {
	return newConversationMsg{}
}

// Add helper function to generate title
func (m MainModel) generateTitle(content string) tea.Cmd {
	return func() tea.Msg {
		titlePrompt := fmt.Sprintf(
			"Please summarize the TEXT below using 3 to 5 words. No talking:\nTEXT:\n%s",
			content,
		)

		titleMsg := api.Message{
			Role:    "user",
			Content: titlePrompt,
		}

		resp, err := m.groqClient.SendMessage(titlePrompt, []api.Message{titleMsg})
		if err != nil {
			return errMsg(err)
		}

		return titleResponseMsg(resp.Content)
	}
}