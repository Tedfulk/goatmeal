package ui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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
	"github.com/tedfulk/goatmeal/services/web/tavily"
	"github.com/tedfulk/goatmeal/ui/theme"
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
	isSearchMode bool
	searchDomains []string
	helpView          *HelpView
	totalCodeBlocks int
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
			a.conversationWindow.SetContent("")
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
			if a.isSearchMode {
				a.isSearchMode = false
				a.input.Reset()
				return a, nil
			}
			if a.currentView == "conversations" {
				// Let the conversation list handle its own escape key
				break
			} else if a.currentView == "glamour" || a.currentView == "username" || a.currentView == "apikeys" || a.currentView == "systemprompts" || a.currentView == "theme" {
				// Let the systemprompts view handle its own escape key for nested views
				if a.currentView == "systemprompts" {
					break
				}
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
				input = strings.TrimPrefix(input, "#")
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
				} else if strings.HasPrefix(input, "m") {
					// Handle message copying to clipboard
					if msgNum, err := strconv.Atoi(strings.TrimPrefix(input, "m")); err == nil {
						for _, m := range a.messages {
							if m.ID == msgNum {
								if err := clipboard.WriteAll(m.Content); err != nil {
									a.statusBar.SetError(fmt.Sprintf("Failed to copy: %v", err))
								} else {
									preview := m.Content
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
				}
				a.input.Reset()
				return a, nil
			}

			// Handle search mode
			if a.isSearchMode {
				if input != "" {
					// Remove the leading "/"
					query := strings.TrimPrefix(input, "/")
					
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
					
					// Create and store user message with domain info
					searchMsg := fmt.Sprintf("🔍 Searching for: %s", query)
					if len(domains) > 0 {
						searchMsg += fmt.Sprintf("\nDomains: %s", strings.Join(domains, ", "))
					}
					userMsg := NewMessage(a.nextMessageID, UserMessage, searchMsg, a.config, a.getNextCodeBlockNumber)
					a.messages = append(a.messages, userMsg)
					a.nextMessageID++
					
					// If this is the first message, generate a title and create conversation in DB
					if len(a.messages) == 1 {
						go a.generateTitle(query)
						a.currentConversationID = uuid.New().String()
					}
					
					a.updateConversationView()
					
					// Clear input and search mode
					a.input.Reset()
					a.isSearchMode = false
					
					// Start spinner before search
					a.statusBar.SetSearchMode(true)
					a.statusBar.SetLoading(true)
					
					// Perform search in goroutine
					go func() {
						defer func() {
							// Stop spinner and reset search mode when done
							a.statusBar.SetLoading(false)
							a.statusBar.SetSearchMode(false)
						}()
						
						tavilyClient := tavily.NewClient(a.config.APIKeys["tavily"])
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
					
					// Return spinner tick command
					return a, a.statusBar.spinner.Tick
				}
			}

			if input != "" {
				userInput := input
				
				userMsg := NewMessage(a.nextMessageID, UserMessage, userInput, a.config, a.getNextCodeBlockNumber)
				a.messages = append(a.messages, userMsg)
				a.nextMessageID++

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
						// Stop spinner when done
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
						for _, msg := range a.messages {
							if msg.Type == UserMessage {
								conversationHistory = append(conversationHistory, "User: "+msg.Content)
							} else if msg.Type == ProviderMessage {
								conversationHistory = append(conversationHistory, "Assistant: "+msg.Content)
							}
						}
						
						// Join history with newlines
						history := strings.Join(conversationHistory, "\n")
						
						// Add current message
						fullPrompt := history + "\nUser: " + userInput

						switch providerName {
						case "anthropic":
							provider := anthropic.NewProvider(apiKey)
							response, err = provider.SendMessage(context.Background(), fullPrompt, a.config.CurrentSystemPrompt, a.config.CurrentModel)
						case "gemini":
							provider := gemini.NewProvider(apiKey)
							response, err = provider.SendMessage(context.Background(), fullPrompt, a.config.CurrentSystemPrompt, a.config.CurrentModel)
						default:
							cfg := providers.OpenAICompatibleConfig{
								Name:    providerName,
								APIKey:  apiKey,
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
							// Save the entire conversation for the first exchange
							conv := &database.Conversation{
								ID:        a.currentConversationID,
								Title:     a.statusBar.conversationTitle,
								Provider:  a.config.CurrentProvider,
								Model:     a.config.CurrentModel,
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
								Messages:  make([]database.Message, len(a.messages)),
							}

							// Convert UI messages to database messages
							for i, msg := range a.messages {
								role := "user"
								if msg.Type == SearchMessage {
									role = "search"
								} else if msg.Type == ProviderMessage {
									role = "assistant"
								}
								conv.Messages[i] = database.Message{
										ID:             uuid.New().String(),
										ConversationID: conv.ID,
										Role:           role,
										Content:        msg.Content,
										CreatedAt:     msg.Timestamp,
									}
							}

							if err := a.db.SaveConversation(conv); err != nil {
								fmt.Printf("Error saving conversation: %v\n", err)
							}

							a.refreshConversationList()
						} else if len(a.messages) > 2 {
							lastMsg := a.messages[len(a.messages)-1]
							role := "user"
							if lastMsg.Type == ProviderMessage {
								role = "assistant"
							}
							dbMsg := &database.Message{
								ID:             uuid.New().String(),
								ConversationID: a.currentConversationID,
								Role:           role,
								Content:        lastMsg.Content,
								CreatedAt:      lastMsg.Timestamp,
							}
							if err := a.db.AddMessage(dbMsg); err != nil {
								fmt.Printf("Error adding message: %v\n", err)
							}
						}
					}
				}()
				
				// Return spinner tick command
				return a, a.statusBar.spinner.Tick
			}
		}

		// Handle menu toggle
		if msg.String() == "?" {
			if a.input.Value() == "" {
				a.showMenu = !a.showMenu
				return a, nil
			}
		}

		// Handle search mode toggle
		if msg.String() == "/" {
			if a.input.Value() == "" {
				a.isSearchMode = true
				a.statusBar.SetSearchMode(true)
				a.input.Set("/")
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

		// In the Update method, where search mode is handled:
		if msg.String() == "backspace" {
			if a.input.Value() == "/" {
				a.isSearchMode = false
				a.statusBar.SetSearchMode(false)
				a.input.Reset()
				return a, nil
			}
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

		// Re-render messages with new width
		a.updateConversationView()
		a.menu.SetSize(msg.Width, msg.Height)

	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
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
		return a.settingsView()
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
