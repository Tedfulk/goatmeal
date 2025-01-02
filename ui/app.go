package ui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/services/providers"
	"github.com/tedfulk/goatmeal/services/providers/anthropic"
	"github.com/tedfulk/goatmeal/services/providers/gemini"
)

var (
	// Styles
	appStyle = lipgloss.NewStyle().
		Padding(0, 1)
	
	// Colors
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#FFF")
	accentColor    = lipgloss.Color("#48BB78")
)

// App represents the main application UI
type App struct {
	config             *config.Config
	conversationWindow viewport.Model
	input             Input
	statusBar         *StatusBar
	height            int
	width             int
	currentView       string
	err               error
	messages          []Message
	nextMessageID     int
	menu             Menu
	showMenu         bool
	settingsMenu     SettingsMenu
	glamourMenu      GlamourMenu
	usernameSettings  UsernameSettings
	apiKeySettings    APIKeySettings
	systemPromptSettings SystemPromptSettings
}

// NewApp creates a new application UI
func NewApp(cfg *config.Config) *App {
	// Initialize viewport with default size, it will be updated in first WindowSizeMsg
	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Align(lipgloss.Left)

	return &App{
		config:            cfg,
		currentView:       "chat",
		input:            NewInput(),
		statusBar:        NewStatusBar(cfg, "New Conversation"),
		messages:         make([]Message, 0),
		nextMessageID:    1,
		conversationWindow: vp,
		menu:            NewMenu(),
		showMenu:        false,
		settingsMenu:    NewSettingsMenu(),
		glamourMenu:     NewGlamourMenu(cfg),
		usernameSettings: NewUsernameSettings(cfg),
		apiKeySettings:   NewAPIKeySettings(cfg),
		systemPromptSettings: NewSystemPromptSettings(cfg),
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return nil
}

// Update handles UI updates
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		case "ctrl+s":
			a.currentView = "settings"
		case "esc":
			if a.currentView == "glamour" || a.currentView == "username" || a.currentView == "apikeys" || a.currentView == "systemprompts" {
				a.currentView = "settings"
			} else if a.currentView == "settings" {
				a.currentView = "chat"
			} else if a.showMenu {
				a.showMenu = false
			}
			return a, nil
		case "enter":
			input := a.input.Value()
			if strings.HasPrefix(input, "#") {
				// Handle message opening
				if msgNum, err := strconv.Atoi(strings.TrimPrefix(input, "#")); err == nil {
					// Find the message
					for _, m := range a.messages {
						if m.ID == msgNum {
							go a.openMessageInEditor(m)
							break
						}
					}
					a.input.Reset()
					return a, nil
				}
			}

			if input != "" {
				userInput := input
				
				// Create and store user message
				userMsg := NewMessage(a.nextMessageID, UserMessage, userInput, a.config)
				a.messages = append(a.messages, userMsg)
				a.nextMessageID++

				// Update conversation window
				a.updateConversationView()
				
				// Clear input and scroll
				a.input.Reset()

				// If this is the first message, generate a title
				if len(a.messages) == 1 {
					go a.generateTitle(userInput)
				}

				// Send message to provider
				go func() {
					var response string
					var err error

					// Get provider instance based on current provider
					providerName := strings.ToLower(a.config.CurrentProvider)
					apiKey := a.config.APIKeys[providerName]

					if apiKey == "" {
						response = fmt.Sprintf("Error: Please provide an API key for %s in the settings", a.config.CurrentProvider)
					} else {
						switch providerName {
						case "anthropic":
							provider := anthropic.NewProvider(apiKey)
							response, err = provider.SendMessage(context.Background(), userInput, a.config.CurrentSystemPrompt, a.config.CurrentModel)
						case "gemini":
							provider := gemini.NewProvider(apiKey)
							response, err = provider.SendMessage(context.Background(), userInput, a.config.CurrentSystemPrompt, a.config.CurrentModel)
						default:
							// Handle OpenAI-compatible providers
							cfg := providers.OpenAICompatibleConfig{
								Name:    providerName,
								APIKey:  apiKey,
							}
							provider := providers.NewOpenAICompatibleProvider(cfg)
							response, err = provider.SendMessage(context.Background(), userInput, a.config.CurrentSystemPrompt, a.config.CurrentModel)
						}

						if err != nil {
							response = "Error: " + err.Error()
						}
					}

					// Create and store provider message
					providerMsg := NewMessage(a.nextMessageID, ProviderMessage, response, a.config)
					a.messages = append(a.messages, providerMsg)
					a.nextMessageID++

					// Update conversation window
					a.updateConversationView()
				}()
			}
		}

		// Handle menu toggle
		if msg.String() == "m" || msg.String() == "ctrl+m" {
			if a.input.Value() == "" {
				a.showMenu = !a.showMenu
				return a, nil
			}
		}

		// If menu is shown, handle its input
		if a.showMenu {
			switch msg.String() {
			case "esc":
				a.showMenu = false
				return a, nil
			case "enter":
				selected := a.menu.list.SelectedItem().(MenuItem)
				a.showMenu = false
				
				// Handle menu selection
				switch selected.title {
				case "New Conversation":
					// TODO: Implement new conversation
				case "Conversations":
					// TODO: Implement conversation list
				case "Settings":
					a.currentView = "settings"
				case "Help":
					// TODO: Implement help
				case "Quit":
					return a, tea.Quit
				}
				
				return a, nil
			}
			
			var menuCmd tea.Cmd
			a.menu, menuCmd = a.menu.Update(msg)
			return a, menuCmd
		}

	case tea.WindowSizeMsg:
		a.height = msg.Height
		a.width = msg.Width
		
		// Update viewport and input sizes
		a.conversationWindow.Width = msg.Width - 4
		a.conversationWindow.Height = msg.Height - 6
		a.input.Width = msg.Width - 4
		a.statusBar.SetWidth(msg.Width)
		a.settingsMenu.SetSize(msg.Width, msg.Height)
		a.glamourMenu.SetSize(msg.Width, msg.Height)
		a.usernameSettings.SetSize(msg.Width, msg.Height)
		a.apiKeySettings.SetSize(msg.Width, msg.Height)
		a.systemPromptSettings.SetSize(msg.Width, msg.Height)

		// Re-render messages with new width
		a.updateConversationView()
		a.menu.SetSize(msg.Width, msg.Height)
	}

	// Update child components
	if a.currentView == "settings" {
		var settingsCmd tea.Cmd
		a.settingsMenu, settingsCmd = a.settingsMenu.Update(msg)
		
		// Check if we need to switch views based on settings menu selection
		if a.settingsMenu.currentView != "settings" {
			a.currentView = a.settingsMenu.currentView
			a.settingsMenu.currentView = "settings" // Reset settings menu view
		}
		
		cmds = append(cmds, settingsCmd)
	} else if a.currentView == "glamour" {
		var glamourCmd tea.Cmd
		a.glamourMenu, glamourCmd = a.glamourMenu.Update(msg)
		cmds = append(cmds, glamourCmd)
	} else if a.currentView == "username" {
		var usernameCmd tea.Cmd
		a.usernameSettings, usernameCmd = a.usernameSettings.Update(msg)
		cmds = append(cmds, usernameCmd)
	} else if a.currentView == "apikeys" {
		var apiKeyCmd tea.Cmd
		a.apiKeySettings, apiKeyCmd = a.apiKeySettings.Update(msg)
		cmds = append(cmds, apiKeyCmd)
	} else if a.currentView == "systemprompts" {
		var systemPromptCmd tea.Cmd
		a.systemPromptSettings, systemPromptCmd = a.systemPromptSettings.Update(msg)
		cmds = append(cmds, systemPromptCmd)
	} else {
		a.conversationWindow, cmd = a.conversationWindow.Update(msg)
		cmds = append(cmds, cmd)

		a.input, cmd = a.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

// updateConversationView updates the conversation window content
func (a *App) updateConversationView() {
	var content string
	for _, msg := range a.messages {
		content += msg.View(a.width) + "\n\n"
	}
	a.conversationWindow.SetContent(content)
}

// View renders the UI
func (a *App) View() string {
	if a.showMenu {
		return a.menu.View()
	}

	switch a.currentView {
	case "settings":
		return a.settingsView()
	case "glamour":
		return a.glamourMenu.View()
	case "username":
		return a.usernameSettings.View()
	case "apikeys":
		return a.apiKeySettings.View()
	case "systemprompts":
		return a.systemPromptSettings.View()
	default:
		return a.chatView()
	}
}

// chatView renders the main chat interface
func (a *App) chatView() string {
	return appStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			a.statusBar.View(),
			a.conversationWindow.View(),
			a.input.View(),
		),
	)
}

// settingsView renders the settings interface
func (a *App) settingsView() string {
	return a.settingsMenu.View()
}

// generateTitle sends a request to generate a conversation title
func (a *App) generateTitle(userInput string) {
	const titleSystemPrompt = `Create a concise, 3-5 word phrase as a header for the following query, strictly adhering to the 3-5 word limit and avoiding the use of the word 'title', and do not generate any other text than the 3-5 word summary and do not use any markdown formatting or any asterisks for bold, and do NOT use quotation marks. 
	Examples of titles:
		Stock Market Trends
		Perfect Chocolate Chip Recipe
		Evolution of Music Streaming
		Remote Work Productivity Tips
		Artificial Intelligence in Healthcare
		Video Game Development Insights`

	apiKey := a.config.APIKeys["groq"]
	if apiKey != "" {
		cfg := providers.OpenAICompatibleConfig{
			Name:    "groq",
			APIKey:  apiKey,
		}
		provider := providers.NewOpenAICompatibleProvider(cfg)
		title, err := provider.SendMessage(context.Background(), userInput, titleSystemPrompt, "llama-3.3-70b-versatile")
		if err == nil {
			a.statusBar.SetConversationTitle(title)
		}
	}
}

// openMessageInEditor opens the message content in the default editor
func (a *App) openMessageInEditor(m Message) {
	// Create a temporary file with read/write permissions
	tmpFile, err := os.CreateTemp("", "goatmeal-*.txt")
	if err != nil {
		return
	}
	tmpPath := tmpFile.Name()

	// Write message content to file
	if _, err := tmpFile.WriteString(m.Content); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return
	}
	tmpFile.Close()

	// Get the default editor from environment variables
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		// Try common editors in order of preference
		if _, err := exec.LookPath("nvim"); err == nil {
			editor = "nvim"
		} else if _, err := exec.LookPath("nano"); err == nil {
			editor = "nano"
		} else {
			editor = "vim"
		}
	}

	cmd := exec.Command(editor, tmpPath)
	if editor != "cursor" {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		os.Remove(tmpPath)
		return
	}

	os.Remove(tmpPath)
} 