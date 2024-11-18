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
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
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
				}
				return m, tea.Quit
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
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

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
		helpText,  // Use our custom help text instead of m.table.HelpView()
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
