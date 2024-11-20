package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

// ModelResponse matches the Groq API response structure
type ModelResponse struct {
	Object string `json:"object"`
	Data   []struct {
		ID string `json:"id"`
	} `json:"data"`
}

type Model struct {
	textInput   textinput.Model
	err         error
	baseStyle   lipgloss.Style
	helpStyle   lipgloss.Style
}

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	terminalWidth     = 80 // default width
	terminalHeight    = 30 // default height
)

type ModelSelectionModel struct {
	table table.Model
}

type item string

func (i item) FilterValue() string { return string(i) }

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter API Key"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 70

	// Create base styles
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("white")).
		Padding(1).
		Width(80).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true)

	// Create help style
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1).
		Align(lipgloss.Center)

	return Model{
		textInput: ti,
		err:       nil,
		baseStyle: baseStyle,
		helpStyle: helpStyle,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		terminalWidth = msg.Width
		terminalHeight = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			apiKey := m.textInput.Value()
			
			// First fetch the models
			models, err := FetchModels(apiKey)
			if err != nil {
				m.err = err
				fmt.Printf("Error fetching models: %v\n", err)
				return m, nil
			}

			// Then save the config
			if err := saveConfig(apiKey); err != nil {
				m.err = err
				fmt.Printf("Error saving config: %v\n", err)
				return m, nil
			}

			// Create and return the model selection screen
			return NewModelSelectionModel(models), nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress any key to quit.", m.err)
	}

	// Create API Key title style
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("white")).
		Padding(0, 1).
		MarginBottom(0).
		Width(80).
		Align(lipgloss.Center)

	// Create the box with input - using direct assignment instead of Copy()
	styledContent := m.baseStyle.
		Render(m.textInput.View())

	// API Key title
	title := titleStyle.Render("API Key")

	// Help text below the box
	helpText := m.helpStyle.Render("Press Enter to submit • Esc to quit")

	// Combine all elements vertically
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		styledContent,
		helpText,
	)

	// Center everything in the terminal
	return lipgloss.Place(
		terminalWidth,
		terminalHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func NewModelSelectionModel(models []string) ModelSelectionModel {
	// Create columns
	columns := []table.Column{
		{Title: "No.", Width: 4},
		{Title: "Model Name", Width: 70},
	}

	// Create rows
	rows := make([]table.Row, len(models))
	for i, model := range models {
		rows[i] = []string{fmt.Sprintf("%d", i+1), model}
	}

	// Initialize table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(13),
	)

	// Style the table
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("white")).
		BorderBottom(true).
		Bold(false).
		Foreground(lipgloss.Color("170"))
	
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("170")).
		Bold(true)
	
	t.SetStyles(s)

	return ModelSelectionModel{table: t}
}

func (m ModelSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			selectedRow := m.table.SelectedRow()
			if len(selectedRow) >= 2 {
				if err := saveSelectedModel(selectedRow[1]); err != nil {
					fmt.Printf("Error saving selected model: %v\n", err)
					return m, tea.Quit
				}
				// Instead of quitting, ask about system prompt
				return NewSystemPromptConfirmModel(), nil
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ModelSelectionModel) View() string {
	// Create title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Padding(0, 0, 1, 0).
		Width(80).
		Align(lipgloss.Center).
		Render("Select a Model")

	// Create base style for the table
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("white"))

	// Create custom help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center)

	helpText := helpStyle.Render("↑/↓: Navigate • enter: Select • q: Quit")

	// Combine elements with proper spacing
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		baseStyle.Render(m.table.View()),
		helpText,
	)

	return lipgloss.Place(
		terminalWidth,
		terminalHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m ModelSelectionModel) Init() tea.Cmd {
	return nil
}

func FetchModels(apiKey string) ([]string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.groq.com/openai/v1/models", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch models: %s", resp.Status)
	}

	var modelResponse ModelResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Filter and extract model IDs
	var models []string
	for _, model := range modelResponse.Data {
		// Skip models containing "whisper"
		if !strings.Contains(strings.ToLower(model.ID), "whisper") {
			models = append(models, model.ID)
		}
	}

	// Sort the models alphabetically
	sort.Strings(models)

	return models, nil
}

func saveConfig(apiKey string) error {
	usr, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %v", err)
	}
	
	configPath := filepath.Join(usr, ".goatmeal", "config.yaml")

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}

	viper.SetConfigFile(configPath)
	viper.Set("api_key", apiKey)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("error writing config: %v", err)
	}

	return nil
}

func saveSelectedModel(model string) error {
	usr, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %v", err)
	}
	
	configPath := filepath.Join(usr, ".goatmeal", "config.yaml")
	
	viper.SetConfigFile(configPath)
	viper.Set("default_model", model)
	
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("error writing config: %v", err)
	}

	return nil
}

// First, add a new model type for the system prompt input
type SystemPromptModel struct {
    textInput textinput.Model
    err       error
}

// Add new model for the yes/no confirmation
type SystemPromptConfirmModel struct {
    quitting bool
}

func NewSystemPromptConfirmModel() SystemPromptConfirmModel {
    return SystemPromptConfirmModel{}
}

func (m SystemPromptConfirmModel) Init() tea.Cmd {
    return nil
}

func (m SystemPromptConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            m.quitting = true
            return m, tea.Quit
        case "n", "N":
            // Save the default system prompt
            if err := saveSystemPrompt(defaultSystemPrompt); err != nil {
                fmt.Printf("Error saving default system prompt: %v\n", err)
            }
            // Transition to theme selection instead of username
            return NewThemeConfirmModel(), nil
        case "y", "Y":
            return NewSystemPromptInputModel(), nil
        }
    }
    return m, nil
}

func (m SystemPromptConfirmModel) View() string {
    if m.quitting {
        return quitTextStyle.Render("Goodbye!")
    }

    // Create confirmation message with center alignment
    text := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("170")).
        Width(78).  // Slightly less than container width to account for padding
        Align(lipgloss.Center).
        Render("Would you like to set a custom system prompt?")

    // Create help text with default prompt info, centered
    defaultInfo := lipgloss.NewStyle().
        Foreground(lipgloss.Color("241")).
        Width(78).  // Slightly less than container width to account for padding
        Align(lipgloss.Center).
        Render("y: Yes • n: No • A default prompt will be used if you select 'n'")


    // Create a container with border
    containerStyle := lipgloss.NewStyle().
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("white")).
        Padding(1).
        Width(80).
        Align(lipgloss.Center)  // Center the container itself

    // Combine elements inside the container
    innerContent := lipgloss.JoinVertical(
        lipgloss.Center,
        text,
        "\n",
        defaultInfo,
    )

    // Wrap in container with border
    content := containerStyle.Render(innerContent)

    // Center in terminal
    return lipgloss.Place(
        terminalWidth,
        terminalHeight,
        lipgloss.Center,
        lipgloss.Center,
        content,
    )
}

// Add new model for the system prompt input
func NewSystemPromptInputModel() SystemPromptModel {
    ti := textinput.New()
    ti.Placeholder = "Enter system prompt"
    ti.Focus()
    ti.Width = 70
    ti.CharLimit = 500

    return SystemPromptModel{
        textInput: ti,
    }
}

func (m SystemPromptModel) Init() tea.Cmd {
    return textinput.Blink
}

func (m SystemPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "esc":
            return m, tea.Quit
        case "enter":
            if err := saveSystemPrompt(m.textInput.Value()); err != nil {
                m.err = err
                return m, tea.Quit
            }
            // Transition to theme selection instead of username
            return NewThemeConfirmModel(), nil
        }
    }

    m.textInput, cmd = m.textInput.Update(msg)
    return m, cmd
}

func (m SystemPromptModel) View() string {
    if m.err != nil {
        return fmt.Sprintf("Error: %v\n\nPress any key to quit.", m.err)
    }

    // Create title
    title := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("170")).
        Padding(0, 0, 1, 0).
        Width(80).
        Align(lipgloss.Center).
        Render("Enter System Prompt")

    // Create input box style
    inputStyle := lipgloss.NewStyle().
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("white")).
        Padding(1).
        Width(80)

    // Create help style with fixed width and center alignment
    helpStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("241")).
        Width(80).  // Added fixed width
        Align(lipgloss.Center)

    // Combine elements
    content := lipgloss.JoinVertical(
        lipgloss.Center,
        title,
        inputStyle.Render(m.textInput.View()),
        helpStyle.Render("Enter to save • Esc to quit"),
    )

    // Center in terminal
    return lipgloss.Place(
        terminalWidth,
        terminalHeight,
        lipgloss.Center,
        lipgloss.Center,
        content,
    )
}

// Add helper function to save the system prompt
func saveSystemPrompt(prompt string) error {
    usr, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("error getting home directory: %v", err)
    }
    
    configPath := filepath.Join(usr, ".goatmeal", "config.yaml")
    
    viper.SetConfigFile(configPath)
    viper.Set("system_prompt", prompt)
    
    if err := viper.WriteConfig(); err != nil {
        return fmt.Errorf("error writing config: %v", err)
    }

    return nil
}

// Keep only one set of constants at the top of the file
const (
    defaultSystemPrompt = "You are a helpful AI assistant. You aim to give accurate, helpful, and concise responses."
    defaultTheme       = "Default"
)

// Keep the rest of the file structure but remove duplicates
type Config struct {
    APIKey        string            `mapstructure:"api_key"`
    DefaultModel  string            `mapstructure:"default_model"`
    SystemPrompt  string            `mapstructure:"system_prompt"`  // Current active prompt
    SystemPrompts []string          `mapstructure:"system_prompts"` // List of saved prompts
    Theme         string            `mapstructure:"theme"`
}

type ThemeColors struct {
    UserBubble       string
    UserText         string
    AssistantBubble  string
    AssistantText    string
    Timestamp        string
}

// Consolidate into a single theme definition
var ThemeMap = map[string]ThemeColors{
    "Default": {
        UserBubble:      "99",  // Light blue
        UserText:        "15",  // White
        AssistantBubble: "62",  // Purple
        AssistantText:  "15",   // White
        Timestamp:      "240",  // Gray
    },
    "Matrix": {
        UserBubble:      "86",      // Light green
        UserText:        "15",      // White
        AssistantBubble: "22",      // Dark green
        AssistantText:  "15",      // White
        Timestamp:      "242",     // Gray
    },
    "Dracula": {
        UserBubble:      "141",     // Purple
        UserText:        "15",      // White
        AssistantBubble: "61",      // Light purple
        AssistantText:  "15",      // White
        Timestamp:      "243",     // Light gray
    },
    "Nord": {
        UserBubble:      "110",     // Blue-green
        UserText:        "15",      // White
        AssistantBubble: "109",     // Light blue
        AssistantText:  "15",      // White
        Timestamp:      "251",     // Light gray
    },
    "Monokai": {
        UserBubble:      "197",     // Pink
        UserText:        "15",      // White
        AssistantBubble: "208",     // Orange
        AssistantText:  "15",      // White
        Timestamp:      "252",     // Light gray
    },
}

// Remove duplicate ThemeConfirmModel.Update method and keep only one implementation
func (m ThemeConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            m.quitting = true
            return m, tea.Quit
        case "n", "N":
            // Save default theme
            if err := saveTheme(defaultTheme); err != nil {
                fmt.Printf("Error saving default theme: %v\n", err)
            }
            m.quitting = true
            return m, tea.Quit
        case "y", "Y":
            return NewThemeSelectionModel(), nil
        }
    }
    return m, nil
}

// Add theme confirmation model
type ThemeConfirmModel struct {
    quitting bool
}

// Add theme selection model
type ThemeSelectionModel struct {
    table table.Model
}

// Add theme confirmation methods
func NewThemeConfirmModel() ThemeConfirmModel {
    return ThemeConfirmModel{}
}

func (m ThemeConfirmModel) Init() tea.Cmd {
    return nil
}

func (m ThemeConfirmModel) View() string {
    if m.quitting {
        return quitTextStyle.Render("Goodbye!")
    }

    text := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("170")).
        Width(78).
        Align(lipgloss.Center).
        Render("Would you like to set a custom theme?")

    defaultInfo := lipgloss.NewStyle().
        Foreground(lipgloss.Color("241")).
        Width(78).
        Align(lipgloss.Center).
        Render("y: Yes • n: No • Default theme will be used if you select 'n'")

    containerStyle := lipgloss.NewStyle().
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("white")).
        Padding(1).
        Width(80).
        Align(lipgloss.Center)

    innerContent := lipgloss.JoinVertical(
        lipgloss.Center,
        text,
        "\n",
        defaultInfo,
    )

    content := containerStyle.Render(innerContent)

    return lipgloss.Place(
        terminalWidth,
        terminalHeight,
        lipgloss.Center,
        lipgloss.Center,
        content,
    )
}

// Add theme selection methods
func NewThemeSelectionModel() ThemeSelectionModel {
    columns := []table.Column{
        {Title: "No.", Width: 4},
        {Title: "Theme Name", Width: 70},
    }

    // Get theme names from the map keys
    var themeNames []string
    for name := range ThemeMap {
        themeNames = append(themeNames, name)
    }
    
    // Sort theme names alphabetically
    sort.Strings(themeNames)

    // Create rows using sorted theme names
    rows := make([]table.Row, len(themeNames))
    for i, name := range themeNames {
        rows[i] = []string{fmt.Sprintf("%d", i+1), name}
    }

    t := table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
        table.WithFocused(true),
        table.WithHeight(7),
    )

    s := table.DefaultStyles()
    s.Header = s.Header.
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("white")).
        BorderBottom(true).
        Bold(false).
        Foreground(lipgloss.Color("170"))
    
    s.Selected = s.Selected.
        Foreground(lipgloss.Color("170")).
        Bold(true)
    
    t.SetStyles(s)

    return ThemeSelectionModel{table: t}
}

func (m ThemeSelectionModel) Init() tea.Cmd {
    return nil
}

func (m ThemeSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "enter":
            selectedRow := m.table.SelectedRow()
            if len(selectedRow) >= 2 {
                if err := saveTheme(selectedRow[1]); err != nil {
                    fmt.Printf("Error saving theme: %v\n", err)
                    return m, tea.Quit
                }
                // Explicitly transition to theme selection
                return NewThemeConfirmModel(), nil
            }
        }
    }

    m.table, cmd = m.table.Update(msg)
    return m, cmd
}

func (m ThemeSelectionModel) View() string {
    title := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("170")).
        Padding(0, 0, 1, 0).
        Width(80).
        Align(lipgloss.Center).
        Render("Select a Theme")

    baseStyle := lipgloss.NewStyle().
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("white"))

    helpStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("241")).
        Width(80).  // Added fixed width
        Align(lipgloss.Center)

    helpText := helpStyle.Render("↑/↓: Navigate • enter: Select • q: Quit")

    content := lipgloss.JoinVertical(
        lipgloss.Center,
        title,
        baseStyle.Render(m.table.View()),
        helpText,
    )

    return lipgloss.Place(
        terminalWidth,
        terminalHeight,
        lipgloss.Center,
        lipgloss.Center,
        content,
    )
}

// Add helper function to save the theme
func saveTheme(theme string) error {
    usr, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("error getting home directory: %v", err)
    }
    
    configPath := filepath.Join(usr, ".goatmeal", "config.yaml")
    
    viper.SetConfigFile(configPath)
    viper.Set("theme", theme)
    
    if err := viper.WriteConfig(); err != nil {
        return fmt.Errorf("error writing config: %v", err)
    }

    return nil
}

// Add this method to get theme colors
func (c *Config) GetThemeColors() ThemeColors {
	if colors, ok := ThemeMap[c.Theme]; ok {
		return colors
	}
	return ThemeMap["Default"]
}

// Add method to save system prompt
func (c *Config) SaveSystemPrompt(prompt string) error {
	usr, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}
	
	configPath := filepath.Join(usr, ".goatmeal", "config.yaml")
	
	viper.SetConfigFile(configPath)
	viper.Set("system_prompt", prompt)
	
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("error writing config: %w", err)
	}

	c.SystemPrompt = prompt
	return nil
}

// Add method to save API key
func (c *Config) SaveAPIKey(apiKey string) error {
    usr, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("error getting home directory: %w", err)
    }
    
    configPath := filepath.Join(usr, ".goatmeal", "config.yaml")
    
    viper.SetConfigFile(configPath)
    viper.Set("api_key", apiKey)
    
    if err := viper.WriteConfig(); err != nil {
        return fmt.Errorf("error writing config: %w", err)
    }

    c.APIKey = apiKey
    return nil
}

// Add method to save theme
func (c *Config) SaveTheme(theme string) error {
    usr, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("error getting home directory: %w", err)
    }
    
    configPath := filepath.Join(usr, ".goatmeal", "config.yaml")
    
    viper.SetConfigFile(configPath)
    viper.Set("theme", theme)
    
    if err := viper.WriteConfig(); err != nil {
        return fmt.Errorf("error writing config: %w", err)
    }

    c.Theme = theme
    return nil
}

// Add this function to handle config loading
func LoadConfig() (*Config, error) {
    usr, err := os.UserHomeDir()
    if err != nil {
        return nil, fmt.Errorf("error getting home directory: %w", err)
    }

    configPath := filepath.Join(usr, ".goatmeal", "config.yaml")
    viper.SetConfigFile(configPath)

    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("error reading config file: %w", err)
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("error unmarshaling config: %w", err)
    }

    // Set defaults if not specified
    if config.DefaultModel == "" {
        config.DefaultModel = "mixtral-8x7b-32768"
    }
    if config.SystemPrompt == "" {
        config.SystemPrompt = defaultSystemPrompt
    }
    if config.Theme == "" {
        config.Theme = defaultTheme
    }

    return &config, nil
}

// Add method to save a new system prompt
func (c *Config) AddSystemPrompt(prompt string) error {
    usr, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("error getting home directory: %w", err)
    }
    
    configPath := filepath.Join(usr, ".goatmeal", "config.yaml")
    
    // Add to prompts list if not already present
    exists := false
    for _, p := range c.SystemPrompts {
        if p == prompt {
            exists = true
            break
        }
    }
    if !exists {
        c.SystemPrompts = append(c.SystemPrompts, prompt)
    }
    
    viper.SetConfigFile(configPath)
    viper.Set("system_prompts", c.SystemPrompts)
    
    if err := viper.WriteConfig(); err != nil {
        return fmt.Errorf("error writing config: %w", err)
    }

    return nil
}

// Add method to set active system prompt
func (c *Config) SetActiveSystemPrompt(prompt string) error {
    usr, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("error getting home directory: %w", err)
    }
    
    configPath := filepath.Join(usr, ".goatmeal", "config.yaml")
    
    viper.SetConfigFile(configPath)
    viper.Set("system_prompt", prompt)
    
    if err := viper.WriteConfig(); err != nil {
        return fmt.Errorf("error writing config: %w", err)
    }

    c.SystemPrompt = prompt
    return nil
}
