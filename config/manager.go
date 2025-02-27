package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// SystemPrompt represents a system prompt with a title and content
type SystemPrompt struct {
	Title   string `mapstructure:"title"`
	Content string `mapstructure:"content"`
}

// Config represents the application configuration
type Config struct {
	APIKeys             map[string]string `mapstructure:"api_keys"`
	CurrentProvider     string           `mapstructure:"current_provider"`
	CurrentModel        string           `mapstructure:"current_model"`
	CurrentSystemPrompt string           `mapstructure:"current_system_prompt"`
	SystemPrompts       []SystemPrompt   `mapstructure:"system_prompts"`
	Settings           Settings         `mapstructure:"settings"`
}

// Settings represents application settings
type Settings struct {
	OutputGlamour          bool       `mapstructure:"outputglamour"`
	ConversationRetention  int        `mapstructure:"conversationretention"`
	Theme                  ThemeConfig `mapstructure:"theme"`
	Username               string     `mapstructure:"username"`
}

// DefaultSettings contains the default values for settings
var DefaultSettings = Settings{
	OutputGlamour:         true,
	ConversationRetention: 30,
	Theme:                DefaultThemeConfig(),
	Username:             "User",
}

// Manager handles configuration loading and saving
type Manager struct {
	config     *Config
	configPath string
	isFirstRun bool
}

// NewManager creates a new configuration manager
func NewManager() (*Manager, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "config.yaml")
	isFirstRun := !fileExists(configPath)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating config directory: %w", err)
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Set default configuration
	setDefaultConfig()

	var config Config
	if !isFirstRun {
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading config: %w", err)
		}

		// Handle theme configuration transition
		var oldTheme string
		if err := viper.UnmarshalKey("settings.theme", &oldTheme); err == nil && oldTheme != "" {
			// We found an old string-based theme config, convert it to the new format
			viper.Set("settings.theme", ThemeConfig{Name: oldTheme})
		}

		if err := viper.Unmarshal(&config); err != nil {
			return nil, fmt.Errorf("error parsing config: %w", err)
		}

		// Ensure theme has a name
		if config.Settings.Theme.Name == "" {
			config.Settings.Theme = DefaultThemeConfig()
		}
	} else {
		// For first run, use the default configuration
		if err := viper.Unmarshal(&config); err != nil {
			return nil, fmt.Errorf("error creating default config: %w", err)
		}
	}

	return &Manager{
		config:     &config,
		configPath: configPath,
		isFirstRun: isFirstRun,
	}, nil
}

// IsFirstRun returns whether this is the first time running the application
func (m *Manager) IsFirstRun() bool {
	return m.isFirstRun
}

// Save writes the current configuration to disk
func (m *Manager) Save() error {
	if err := viper.WriteConfig(); err != nil {
		if os.IsNotExist(err) {
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("error creating config file: %w", err)
			}
		} else {
			return fmt.Errorf("error saving config: %w", err)
		}
	}
	return nil
}

// GetConfig returns the current configuration
func (m *Manager) GetConfig() *Config {
	return m.config
}

// SetAPIKey updates an API key in the configuration
func (m *Manager) SetAPIKey(provider, key string) error {
	if m.config.APIKeys == nil {
		m.config.APIKeys = make(map[string]string)
	}
	m.config.APIKeys[provider] = key
	viper.Set("api_keys", m.config.APIKeys)
	return m.Save()
}

// GetAPIKey retrieves an API key from the configuration
func (m *Manager) GetAPIKey(provider string) string {
	if m.config.APIKeys == nil {
		return ""
	}
	return m.config.APIKeys[provider]
}

// SetCurrentProvider sets the current provider
func (m *Manager) SetCurrentProvider(provider string) error {
	m.config.CurrentProvider = provider
	viper.Set("current_provider", provider)
	return m.Save()
}

// SetCurrentModel sets the current model
func (m *Manager) SetCurrentModel(model string) error {
	m.config.CurrentModel = model
	viper.Set("current_model", model)
	return m.Save()
}

// SetCurrentSystemPrompt sets the current system prompt
func (m *Manager) SetCurrentSystemPrompt(prompt string) error {
	m.config.CurrentSystemPrompt = prompt
	viper.Set("current_system_prompt", prompt)
	return m.Save()
}

// AddSystemPrompt adds a new system prompt to the configuration
func (m *Manager) AddSystemPrompt(title, content string) error {
	prompt := SystemPrompt{
		Title:   title,
		Content: content,
	}
	m.config.SystemPrompts = append(m.config.SystemPrompts, prompt)
	viper.Set("system_prompts", m.config.SystemPrompts)
	return m.Save()
}

// DeleteSystemPrompt removes a system prompt from the configuration
func (m *Manager) DeleteSystemPrompt(title string) error {
	prompts := make([]SystemPrompt, 0)
	for _, p := range m.config.SystemPrompts {
		if p.Title != title {
			prompts = append(prompts, p)
		}
	}
	m.config.SystemPrompts = prompts
	viper.Set("system_prompts", m.config.SystemPrompts)
	return m.Save()
}

// UpdateSettings updates the application settings
func (m *Manager) UpdateSettings(settings Settings) error {
	m.config.Settings = settings
	viper.Set("settings", settings)
	return m.Save()
}

// UpdateUsername updates the username in the config file
func (m *Manager) UpdateUsername(username string) error {
	viper.Set("username", username)
	return viper.WriteConfig()
}

// SetSystemPrompts updates the system prompts in the config file
func (m *Manager) SetSystemPrompts(prompts []SystemPrompt) error {
	viper.Set("system_prompts", prompts)
	return viper.WriteConfig()
}

// setDefaultConfig sets the default configuration values
func setDefaultConfig() {
	viper.SetDefault("api_keys", map[string]string{
		"groq":     "",
		"openai":   "",
		"anthropic": "",
		"gemini":   "",
		"deepseek": "",
		"tavily":   "",
		"ollama":   "",
	})

	// Default current selections
	viper.SetDefault("current_provider", "")
	viper.SetDefault("current_model", "")
	viper.SetDefault("current_system_prompt", "")

	// Default system prompts
	viper.SetDefault("system_prompts", []SystemPrompt{
		{
			Title: "General",
			Content: "You are a helpful AI assistant. You aim to give accurate, helpful, and concise responses.",
		},
	})

	// Default settings
	viper.SetDefault("settings", DefaultSettings)
}

// getConfigDir returns the path to the configuration directory
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "goatmeal"), nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Load loads the configuration from disk
func Load() (*Config, error) {
	manager, err := NewManager()
	if err != nil {
		return nil, err
	}
	return manager.GetConfig(), nil
} 