package config

// ThemeConfig represents the theme configuration
type ThemeConfig struct {
	Name string `yaml:"name"`
}

// DefaultThemeConfig returns the default theme configuration
func DefaultThemeConfig() ThemeConfig {
	return ThemeConfig{
		Name: "Default",
	}
} 