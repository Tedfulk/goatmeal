package ui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/database"
	"github.com/tedfulk/goatmeal/services/providers"
	"github.com/tedfulk/goatmeal/services/providers/anthropic"
	"github.com/tedfulk/goatmeal/services/providers/gemini"
	"github.com/tedfulk/goatmeal/services/providers/ollama"
	"github.com/tedfulk/goatmeal/services/search"
	"github.com/tedfulk/goatmeal/ui/theme"
	"github.com/tedfulk/goatmeal/utils/editor"
	"github.com/tedfulk/goatmeal/utils/prompts"
)

var (
	appStyle = lipgloss.NewStyle().
		Padding(0, 1)
)

type App struct {
	config             *config.Config
	db                 *database.DB
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
	conversationList  *ConversationListView
	currentConversationID string
	searchDomains []string
	helpView          *HelpView
	totalCodeBlocks int
	queryEnhancer *search.QueryEnhancer
}

func NewApp(cfg *config.Config, db *database.DB) *App {
	theme.LoadThemeFromConfig(cfg.Settings.Theme.Name)

	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.CurrentTheme.Primary.GetColor()).
		Align(lipgloss.Left)

	return &App{
		config:            cfg,
		db:               db,
		currentView:       "chat",
		input:            NewInput(),
		statusBar:        NewStatusBar(cfg, "New Conversation"),
		messages:         make([]Message, 0),
		nextMessageID:    1,
		conversationWindow: vp,
		menu:            NewMenu(),
		showMenu:        false,
		settingsMenu:    NewSettingsMenu(cfg),
		glamourMenu:     NewGlamourMenu(cfg),
		usernameSettings: NewUsernameSettings(cfg),
		apiKeySettings:   NewAPIKeySettings(cfg),
		
		systemPromptSettings: NewSystemPromptSettings(cfg),
		conversationList: NewConversationListView(db, cfg),
		helpView:          NewHelpView(),
		totalCodeBlocks: 0,
		queryEnhancer: search.NewQueryEnhancer(cfg.APIKeys["groq"]),
	}
}

func (a *App) Init() tea.Cmd {
	return tea.EnableMouseCellMotion
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case ModelChangeMsg:
		// Update the config
		a.config.CurrentProvider = msg.Provider
		a.config.CurrentModel = msg.Model
		
		// Update the status bar
		a.statusBar.UpdateProviderAndModel(msg.Provider, msg.Model)
		
		// Return to settings view
		a.currentView = "settings"
		return a, nil

	case ThemeChangeMsg:
		a.conversationWindow.Style = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(theme.CurrentTheme.Primary.GetColor()).
				Align(lipgloss.Left)
		a.statusBar.UpdateStyle()
		a.menu.list.Styles.Title = theme.BaseStyle.Title.
				Foreground(theme.CurrentTheme.Primary.GetColor())
		a.settingsMenu.list.Styles.Title = theme.BaseStyle.Title.
				Foreground(theme.CurrentTheme.Primary.GetColor())
		a.glamourMenu.list.Styles.Title = theme.BaseStyle.Title.
				Foreground(theme.CurrentTheme.Primary.GetColor())
		a.systemPromptSettings.list.Styles.Title = theme.BaseStyle.Title.
				Foreground(theme.CurrentTheme.Primary.GetColor())
		return a, nil
	case SetViewMsg:
		a.currentView = msg.view
		if msg.view == "chat" {
			a.showMenu = true
		} else if msg.view == "conversations" {
			a.refreshConversationList()
		}
		return a, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Batch(
				tea.DisableMouse,
				tea.Quit,
			)
		case "ctrl+s":
			a.currentView = "settings"
		case "ctrl+l":
			a.currentView = "conversations"
			a.refreshConversationList()
		case "ctrl+t":
			// Start a new conversation
			a.messages = make([]Message, 0)
			a.nextMessageID = 1
			a.currentConversationID = ""
			a.totalCodeBlocks = 0  // Reset code block counter
			a.statusBar.SetConversationTitle("New Conversation")
			a.currentView = "chat"
			a.showMenu = false
			a.updateConversationView()
			a.refreshConversationList()
		case "esc":
			if a.currentView == "help" {
				a.currentView = "chat"
				return a, nil
			}
			if a.currentView == "conversations" {
				// Let the conversation list handle its own escape key
				break
			} else if a.currentView == "glamour" || a.currentView == "username" || a.currentView == "apikeys" || a.currentView == "systemprompts" || a.currentView == "theme" || a.currentView == "model" {
				// Let these views handle their own escape key
				break
			} else if a.currentView == "settings" {
				a.currentView = "chat"
			} else if a.showMenu {
				a.showMenu = false
			}
			return a, nil
		case "enter":
			input := a.input.Value()
			if strings.HasPrefix(input, "/") {
				input = strings.TrimPrefix(input, "/")
				if strings.HasPrefix(input, "o") {
					// Handle message opening to default editor
					if msgNum, err := strconv.Atoi(strings.TrimPrefix(input, "o")); err == nil {
						for _, m := range a.messages {
							if m.ID == msgNum {
								go a.openMessageInEditor(m)
								break
							}
						}
					}
				} else if strings.HasPrefix(input, "c") {
					// Handle message copying to clipboard
					if msgNum, err := strconv.Atoi(strings.TrimPrefix(input, "c")); err == nil {
						for _, m := range a.messages {
							if m.ID == msgNum {
								// Get message content and strip search prefixes if present
								contentToCopy := search.StripPrefix(m.Content)
								
								if err := clipboard.WriteAll(contentToCopy); err != nil {
									a.statusBar.SetError(fmt.Sprintf("Failed to copy: %v", err))
								} else {
									preview := contentToCopy
									if len(preview) > 20 {
										preview = preview[:20] + "..."
									}
									a.statusBar.SetTemporaryText(fmt.Sprintf("📋 Copied: %s", preview))
								}
								break
							}
						}
					}
				} else if strings.HasPrefix(input, "b") {
					// Handle code block copying
					if blockNum, err := strconv.Atoi(strings.TrimPrefix(input, "b")); err == nil {
						// Find the last provider message
						for i := len(a.messages) - 1; i >= 0; i-- {
							m := a.messages[i]
							if m.Type == ProviderMessage {
								// Extract code blocks and find the requested one
								blocks := m.ExtractCodeBlocks()
								for _, block := range blocks {
									if block.Number == blockNum {
										if err := clipboard.WriteAll(block.Content); err != nil {
											a.statusBar.SetError(fmt.Sprintf("Failed to copy code block: %v", err))
										} else {
											preview := block.Content
											if len(preview) > 20 {
												preview = preview[:20] + "..."
											}
											a.statusBar.SetTemporaryText(fmt.Sprintf("📋 Copied code block [^%d]: %s", blockNum, preview))
										}
										break
									}
								}
								break
							}
						}
					}
				} else if strings.HasPrefix(input, "s") {
					// Handle message speaking
					if msgNum, err := strconv.Atoi(strings.TrimPrefix(input, "s")); err == nil {
						for _, m := range a.messages {
							if m.ID == msgNum {
								go func(msg Message) {
									if err := msg.Speak(); err != nil {
										a.statusBar.SetError(fmt.Sprintf("Failed to speak message: %v", err))
									}
								}(m)
								break
							}
						}
					}
				} else if strings.HasPrefix(input, "web") || strings.HasPrefix(input, "webe") {
					// Handle web search
					isEnhanced := strings.HasPrefix(input, "webe")
					var query string
					if isEnhanced {
						query = strings.TrimSpace(strings.TrimPrefix(input, "webe"))
					} else {
						query = strings.TrimSpace(strings.TrimPrefix(input, "web"))
					}

					if query != "" {
						// Check for domain inclusions (marked with +)
						var domains []string
						parts := strings.Split(query, "+")
						query = strings.TrimSpace(parts[0])
						
						// If there are additional parts, they are domains
						if len(parts) > 1 {
							for _, domain := range parts[1:] {
								domain = strings.TrimSpace(domain)
								if domain != "" {
									domains = append(domains, domain)
								}
							}
						}

						// Create search message prefix based on mode
						searchMsg := fmt.Sprintf("🔍 Searching for: %s", query)
						if isEnhanced {
							enhanced, err := a.enhanceSearchQuery(query)
							if err != nil {
								a.statusBar.SetError(fmt.Sprintf("Failed to enhance query: %v", err))
								a.input.Reset()
								return a, nil
							}
							query = enhanced
							searchMsg = fmt.Sprintf("🔍+ Enhanced search: %s", query)
						}
						if len(domains) > 0 {
							searchMsg += fmt.Sprintf("\nDomains: %s", strings.Join(domains, ", "))
						}

						// Create and store user message
						userMsg := NewMessage(a.nextMessageID, UserMessage, searchMsg, a.config, a.getNextCodeBlockNumber)
						a.messages = append(a.messages, userMsg)
						a.nextMessageID++

						// If this is the first message, generate a title and create conversation
						if len(a.messages) == 1 {
							go a.generateTitle(query)
							a.currentConversationID = uuid.New().String()
						}

						a.updateConversationView()
						a.input.Reset()

						// Start spinner
						a.statusBar.SetLoading(true)

						// Perform search in goroutine
						go func() {
							defer func() {
								a.statusBar.SetLoading(false)
							}()

							tavilyClient := search.NewClient(a.config.APIKeys["tavily"])
							searchResp, err := tavilyClient.Search(query, domains)
							
							var response string
							if err != nil {
								response = fmt.Sprintf("Error performing search: %v", err)
							} else {
								// Format search results
								var sb strings.Builder
								
								// If there's an answer, show it first
								if searchResp.Answer != nil {
									sb.WriteString("### Answer Summary\n\n")
									sb.WriteString(fmt.Sprintf("%v", searchResp.Answer))
									sb.WriteString("\n\n---\n\n")
								}
								
								sb.WriteString("### Search Results\n\n")
								
								// Show individual results
								for _, result := range searchResp.Results {
									sb.WriteString(fmt.Sprintf("**[%s](%s)**\n\n%s\n\n---\n\n", 
										result.Title, result.URL, result.Content))
								}
								response = sb.String()
							}
							
							// Create and store search result message
							searchMsg := NewMessage(a.nextMessageID, SearchMessage, response, a.config, a.getNextCodeBlockNumber)
							a.messages = append(a.messages, searchMsg)
							a.nextMessageID++
							
							// Save to database if we have a current conversation
							if a.currentConversationID != "" {
								if len(a.messages) == 2 { // First user message + first search response
									// Save the entire conversation
									conv := &database.Conversation{
										ID:        a.currentConversationID,
										Title:     a.statusBar.conversationTitle,
										Provider:  "tavily",
										Model:     "search",
										CreatedAt: time.Now(),
										UpdatedAt: time.Now(),
										Messages:  make([]database.Message, len(a.messages)),
									}

									// Convert UI messages to database messages
									for i, msg := range a.messages {
										role := "user"
										if msg.Type == SearchMessage {
											role = "search"
										}
										conv.Messages[i] = database.Message{
											ID:             uuid.New().String(),
											ConversationID: conv.ID,
											Role:           role,
											Content:        msg.Content,
											CreatedAt:      msg.Timestamp,
										}
									}

									if err := a.db.SaveConversation(conv); err != nil {
										fmt.Printf("Error saving conversation: %v\n", err)
									}
									a.refreshConversationList()
								} else {
									// Add just the new message for subsequent messages
									dbMsg := &database.Message{
										ID:             uuid.New().String(),
										ConversationID: a.currentConversationID,
										Role:           "search",
										Content:        response,
										CreatedAt:      time.Now(),
									}
									if err := a.db.AddMessage(dbMsg); err != nil {
										fmt.Printf("Error adding search message: %v\n", err)
									}
								}
							}
							
							a.updateConversationView()
						}()

						return a, a.statusBar.spinner.Tick
					}
				}
				a.input.Reset()
				return a, nil
			}

			if input != "" {
				userInput := input
				
				// Create and store user message
				userMsg := NewMessage(a.nextMessageID, UserMessage, userInput, a.config, a.getNextCodeBlockNumber)
				a.messages = append(a.messages, userMsg)
				a.nextMessageID++

				// Update view and reset input immediately
				a.updateConversationView()
				a.input.Reset()

				// If this is the first message, generate a title and create conversation in DB
				if len(a.messages) == 1 {
					go a.generateTitle(userInput)
					a.currentConversationID = uuid.New().String()
				}

				// Start spinner before sending message
				a.statusBar.SetLoading(true)
				
				// Send message to provider
				go func() {
					defer func() {
						a.statusBar.SetLoading(false)
					}()

					var response string
					var err error

					// Get provider instance based on current provider
					providerName := strings.ToLower(a.config.CurrentProvider)
					apiKey := a.config.APIKeys[providerName]

					if apiKey == "" {
						response = fmt.Sprintf("Error: Please provide an API key for %s in the settings", a.config.CurrentProvider)
					} else {
						// Build conversation history
						var conversationHistory []string
						for _, msg := range a.messages[:len(a.messages)-1] { // Exclude the message we just added
							if msg.Type == UserMessage {
								conversationHistory = append(conversationHistory, "User: "+msg.Content)
							} else if msg.Type == ProviderMessage {
								conversationHistory = append(conversationHistory, "Assistant: "+msg.Content)
							}
							// Skip SearchMessage type when building conversation history
						}
						
						// Add current message
						conversationHistory = append(conversationHistory, "User: "+userInput)
						
						// Join history with newlines
						fullPrompt := strings.Join(conversationHistory, "\n")

						switch providerName {
						case "anthropic":
							provider := anthropic.NewProvider(apiKey)
							response, err = provider.SendMessage(context.Background(), fullPrompt, a.config.CurrentSystemPrompt, a.config.CurrentModel)
						case "gemini":
							provider := gemini.NewProvider(apiKey)
							response, err = provider.SendMessage(context.Background(), fullPrompt, a.config.CurrentSystemPrompt, a.config.CurrentModel)
						case "ollama":
							provider := ollama.NewProvider(apiKey)
							response, err = provider.SendMessage(context.Background(), fullPrompt, a.config.CurrentSystemPrompt, a.config.CurrentModel)
						default:
							cfg := providers.OpenAICompatibleConfig{
								Name:   providerName,
								APIKey: apiKey,
							}
							provider := providers.NewOpenAICompatibleProvider(cfg)
							response, err = provider.SendMessage(context.Background(), fullPrompt, a.config.CurrentSystemPrompt, a.config.CurrentModel)
						}

						if err != nil {
							response = "Error: " + err.Error()
						}
					}

					// Create and store provider message
					providerMsg := NewMessage(a.nextMessageID, ProviderMessage, response, a.config, a.getNextCodeBlockNumber)
					a.messages = append(a.messages, providerMsg)
					a.nextMessageID++

					a.updateConversationView()

					// Save to database if we have a current conversation
					if a.currentConversationID != "" {
						if len(a.messages) == 2 { // First user message + first AI response
							// Save the entire conversation
							conv := &database.Conversation{
								ID:        a.currentConversationID,
								Title:     a.statusBar.conversationTitle,
								Provider:  a.config.CurrentProvider,
								Model:     a.config.CurrentModel,
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
								Messages:  make([]database.Message, 2), // Only save the current exchange
							}

							// Convert only the current exchange messages
							for i := 0; i < 2; i++ {
								msg := a.messages[i]
								role := "user"
								if msg.Type == ProviderMessage {
									role = "assistant"
								} else if msg.Type == SearchMessage {
									role = "search"
								}
								conv.Messages[i] = database.Message{
									ID:             uuid.New().String(),
									ConversationID: conv.ID,
									Role:           role,
									Content:        msg.Content,
									CreatedAt:      msg.Timestamp,
								}
							}

							if err := a.db.SaveConversation(conv); err != nil {
								fmt.Printf("Error saving conversation: %v\n", err)
							}
							a.refreshConversationList()
						} else {
							// Add just the last two messages (current exchange)
							lastMsgIndex := len(a.messages) - 1
							messages := []database.Message{
								{
									ID:             uuid.New().String(),
									ConversationID: a.currentConversationID,
									Role:           "user",
									Content:        a.messages[lastMsgIndex-1].Content,
									CreatedAt:      a.messages[lastMsgIndex-1].Timestamp,
								},
								{
									ID:             uuid.New().String(),
									ConversationID: a.currentConversationID,
									Role:           "assistant",
									Content:        a.messages[lastMsgIndex].Content,
									CreatedAt:      a.messages[lastMsgIndex].Timestamp,
								},
							}
							
							for _, dbMsg := range messages {
								if err := a.db.AddMessage(&dbMsg); err != nil {
									fmt.Printf("Error adding message: %v\n", err)
								}
							}
						}
					}
				}()
				
				return a, a.statusBar.spinner.Tick
			}
		case "ctrl+q":
			StopSpeech()
			return a, nil
		}

		// Handle menu toggle
		if msg.String() == "?" {
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
				
				switch selected.title {
				case "New Conversation":
					a.messages = make([]Message, 0)
					a.nextMessageID = 1
					a.currentConversationID = ""
					a.statusBar.SetConversationTitle("New Conversation")
					a.updateConversationView()
					a.refreshConversationList()
				case "Conversations":
					a.currentView = "conversations"
					a.refreshConversationList()
				case "Settings":
					a.currentView = "settings"
				case "Help":
					a.currentView = "help"
				case "Quit":
					return a, tea.Quit
				}
				
				return a, nil
			}
			
			var menuCmd tea.Cmd
			a.menu, menuCmd = a.menu.Update(msg)
			return a, menuCmd
		}

		// Handle help view
		if msg.String() == "ctrl+h" {
			a.currentView = "help"
			return a, nil
		}

	case tea.WindowSizeMsg:
		a.height = msg.Height
		a.width = msg.Width
		

		a.conversationWindow.Width = msg.Width - 4
		a.conversationWindow.Height = msg.Height - 6
		a.input.Width = msg.Width - 4
		a.statusBar.SetWidth(msg.Width)
		a.settingsMenu.SetSize(msg.Width, msg.Height)
		a.glamourMenu.SetSize(msg.Width, msg.Height)
		a.usernameSettings.SetSize(msg.Width, msg.Height)
		a.apiKeySettings.SetSize(msg.Width, msg.Height)
		a.systemPromptSettings.SetSize(msg.Width, msg.Height)
		a.conversationList.SetSize(msg.Width, msg.Height)
		a.helpView.SetSize(msg.Width, msg.Height)

		a.updateConversationView()
		a.menu.SetSize(msg.Width, msg.Height)

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
			// Check if click is in status bar
			if a.statusBar.inBounds(msg.X, msg.Y) {
				// Start new conversation
				a.messages = make([]Message, 0)
				a.nextMessageID = 1
				a.currentConversationID = ""
				a.statusBar.SetConversationTitle("New Conversation")
				a.currentView = "chat"
				a.showMenu = false
				a.updateConversationView()
				a.refreshConversationList()
				return a, nil
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		if a.statusBar.isLoading {
			cmd = a.statusBar.Update(msg)
		}
		return a, cmd

	case SystemPromptChangeMsg:
		// Update the app's config with the new system prompt
		a.config.CurrentSystemPrompt = msg.NewPrompt
		// Reload the config to ensure everything is in sync
		if newConfig, err := config.Load(); err == nil {
			a.config = newConfig
			// Update any components that need the new config
			a.statusBar.config = newConfig
		}
		return a, nil
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
	} else if a.currentView == "conversations" {
		var listCmd tea.Cmd
		a.conversationList, listCmd = a.conversationList.Update(msg)
		cmds = append(cmds, listCmd)
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
	} else if a.currentView == "theme" {
		var themeCmd tea.Cmd
		a.settingsMenu.themeSettings, themeCmd = a.settingsMenu.themeSettings.Update(msg)
		cmds = append(cmds, themeCmd)
	} else if a.currentView == "model" {
		var modelCmd tea.Cmd
		a.settingsMenu.modelSettings, modelCmd = a.settingsMenu.modelSettings.Update(msg)
		cmds = append(cmds, modelCmd)
	} else if a.currentView == "help" {
		var helpCmd tea.Cmd
		a.helpView, helpCmd = a.helpView.Update(msg)
		cmds = append(cmds, helpCmd)
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
	var lineCount int
	
	// First pass: count lines up to the last provider message
	for i, msg := range a.messages {
		if i < len(a.messages)-1 {
			content += msg.View(a.width) + "\n\n"
			// Count newlines in the message plus the two we add
			lineCount += strings.Count(msg.View(a.width), "\n") + 2
		}
	}
	
	// Add the last message
	if len(a.messages) > 0 {
		lastMsg := a.messages[len(a.messages)-1]
		content += lastMsg.View(a.width) + "\n\n"
		
		// If it's a provider message or search message, scroll to its position
		if lastMsg.Type == ProviderMessage || lastMsg.Type == SearchMessage {
			a.conversationWindow.SetContent(content)
			a.conversationWindow.GotoTop() // First go to top
			// Then scroll down to the last provider/search message's position
			a.conversationWindow.LineDown(lineCount)
		} else {
			// For user messages, just scroll to bottom as before
			a.conversationWindow.SetContent(content)
			a.conversationWindow.GotoBottom()
		}
		return
	}
	
	// If no messages, just set content
	a.conversationWindow.SetContent(content)
}

// View renders the UI
func (a *App) View() string {
	if a.showMenu {
		return a.menu.View()
	}

	switch a.currentView {
	case "settings":
		return a.settingsMenu.View()
	case "conversations":
		return a.conversationList.View()
	case "glamour":
		return a.glamourMenu.View()
	case "username":
		return a.usernameSettings.View()
	case "apikeys":
		return a.apiKeySettings.View()
	case "systemprompts":
		return a.systemPromptSettings.View()
	case "theme":
		return a.settingsMenu.themeSettings.View()
	case "model":
		return a.settingsMenu.modelSettings.View()
	case "help":
		return a.helpView.View()
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
	apiKey := a.config.APIKeys["groq"]
	if apiKey != "" {
		cfg := providers.OpenAICompatibleConfig{
			Name:    "groq",
			APIKey:  apiKey,
		}
		provider := providers.NewOpenAICompatibleProvider(cfg)
		title, err := provider.SendMessage(context.Background(), userInput, prompts.GetTitleSystemPrompt(), "llama-3.3-70b-versatile")
		if err == nil {
			title = strings.Join(strings.Fields(strings.ReplaceAll(title, "\n", " ")), " ")
			if len(title) > 27 {
				lastSpace := strings.LastIndex(title[:27], " ")
				if lastSpace == -1 {
					title = title[:27]
				} else {
					title = title[:lastSpace]
				}
			}
			
			a.statusBar.SetConversationTitle(title)
			// Update conversation title in database if we have a current conversation
			if a.currentConversationID != "" {
				if err := a.db.UpdateConversationTitle(a.currentConversationID, title); err != nil {
					fmt.Printf("Error updating conversation title: %v\n", err)
				}
			}
		}
	}
}

// openMessageInEditor opens the message content in the default editor
func (a *App) openMessageInEditor(m Message) {
	editor.OpenInEditor(m.Content)
}

// In the App struct, add a method to refresh the conversation list
func (a *App) refreshConversationList() {
	if a.conversationList != nil {
		a.conversationList.loadConversations()
	}
}

func (a *App) getNextCodeBlockNumber() int {
	a.totalCodeBlocks++
	return a.totalCodeBlocks
}

// Add this method to the App struct
func (a *App) enhanceSearchQuery(query string) (string, error) {
	return a.queryEnhancer.Enhance(query)
}
