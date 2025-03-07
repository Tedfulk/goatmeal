package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name    string
	Primary CompleteAdaptiveColor
	Secondary CompleteAdaptiveColor
	Text CompleteAdaptiveColor
	Accent CompleteAdaptiveColor
	StatusBar StatusBarColors
	Border BorderColors
	Message MessageColors
}

type StatusBarColors struct {
	Text       CompleteAdaptiveColor
	Title      CompleteAdaptiveColor
	Model      CompleteAdaptiveColor
}

type BorderColors struct {
	Normal CompleteAdaptiveColor
	Active CompleteAdaptiveColor
}

type MessageColors struct {
	UserText       CompleteAdaptiveColor
	AIText        CompleteAdaptiveColor
	Timestamp   CompleteAdaptiveColor
}

type CompleteAdaptiveColor struct {
	Light CompleteColor
	Dark  CompleteColor
}

type CompleteColor struct {
	TrueColor string
	ANSI256   string
	ANSI      string
}

// BaseStyle contains the base styles for components without colors
var BaseStyle = struct {
	Menu      lipgloss.Style
	Title     lipgloss.Style
	Input     lipgloss.Style
	StatusBar lipgloss.Style
	Message   lipgloss.Style
}{
	Menu: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(50),

	Title: lipgloss.NewStyle().
		Bold(true).
		Width(46).
		Align(lipgloss.Center),

	Input: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()),

	StatusBar: lipgloss.NewStyle(),

	Message: lipgloss.NewStyle().
		Padding(1).
		BorderStyle(lipgloss.RoundedBorder()),
}

// DefaultTheme is the default color theme
var DefaultTheme = Theme{
	Name: "Default",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#7D56F4", ANSI256: "99", ANSI: "5"},
		Dark:  CompleteColor{TrueColor: "#7D56F4", ANSI256: "99", ANSI: "5"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
			Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#6D46E4", ANSI256: "98", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#6D46E4", ANSI256: "98", ANSI: "5"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#9E7DFF", ANSI256: "141", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#9E7DFF", ANSI256: "141", ANSI: "13"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#7D56F4", ANSI256: "99", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#7D56F4", ANSI256: "99", ANSI: "5"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#7D56F4", ANSI256: "99", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#7D56F4", ANSI256: "99", ANSI: "5"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#666666", ANSI256: "241", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#666666", ANSI256: "241", ANSI: "8"},
		},
	},
}

// DraculaTheme is a dark theme with vibrant colors
var DraculaTheme = Theme{
	Name: "Dracula",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#BD93F9", ANSI256: "141", ANSI: "13"},
		Dark:  CompleteColor{TrueColor: "#BD93F9", ANSI256: "141", ANSI: "13"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#F8F8F2", ANSI256: "231", ANSI: "15"},
		Dark:  CompleteColor{TrueColor: "#F8F8F2", ANSI256: "231", ANSI: "15"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#F8F8F2", ANSI256: "231", ANSI: "15"},
			Dark:  CompleteColor{TrueColor: "#F8F8F2", ANSI256: "231", ANSI: "15"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#50FA7B", ANSI256: "84", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#50FA7B", ANSI256: "84", ANSI: "2"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#8BE9FD", ANSI256: "117", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#8BE9FD", ANSI256: "117", ANSI: "6"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#9D73F9", ANSI256: "140", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#9D73F9", ANSI256: "140", ANSI: "13"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF79C6", ANSI256: "212", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#FF79C6", ANSI256: "212", ANSI: "5"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#BD93F9", ANSI256: "141", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#BD93F9", ANSI256: "141", ANSI: "13"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#F8F8F2", ANSI256: "231", ANSI: "15"},
			Dark:  CompleteColor{TrueColor: "#F8F8F2", ANSI256: "231", ANSI: "15"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#6272A4", ANSI256: "61", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#6272A4", ANSI256: "61", ANSI: "8"},
		},
	},
}

// NordTheme is a cool
var NordTheme = Theme{
	Name: "Nord",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#88C0D0", ANSI256: "110", ANSI: "4"},
		Dark:  CompleteColor{TrueColor: "#88C0D0", ANSI256: "110", ANSI: "4"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#ECEFF4", ANSI256: "231", ANSI: "15"},
		Dark:  CompleteColor{TrueColor: "#ECEFF4", ANSI256: "231", ANSI: "15"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#ECEFF4", ANSI256: "231", ANSI: "15"},
			Dark:  CompleteColor{TrueColor: "#ECEFF4", ANSI256: "231", ANSI: "15"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#A3BE8C", ANSI256: "144", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#A3BE8C", ANSI256: "144", ANSI: "2"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#88C0D0", ANSI256: "110", ANSI: "4"},
			Dark:  CompleteColor{TrueColor: "#88C0D0", ANSI256: "110", ANSI: "4"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#81B5C5", ANSI256: "109", ANSI: "4"},
			Dark:  CompleteColor{TrueColor: "#81B5C5", ANSI256: "109", ANSI: "4"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#B48EAD", ANSI256: "139", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#B48EAD", ANSI256: "139", ANSI: "5"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#88C0D0", ANSI256: "110", ANSI: "4"},
			Dark:  CompleteColor{TrueColor: "#88C0D0", ANSI256: "110", ANSI: "4"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#88C0D0", ANSI256: "110", ANSI: "4"},
			Dark:  CompleteColor{TrueColor: "#88C0D0", ANSI256: "110", ANSI: "4"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#4C566A", ANSI256: "59", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#4C566A", ANSI256: "59", ANSI: "8"},
		},
	},
}

func (c CompleteAdaptiveColor) GetColor() lipgloss.Color {
	// For now, we'll just use the dark variant's true color
	// TODO: Implement proper light/dark detection and color profile handling
	return lipgloss.Color(c.Dark.TrueColor)
}

// CurrentTheme holds the currently active theme
var CurrentTheme = DefaultTheme 

// MatrixClassicTheme is the classic green matrix theme
var MatrixClassicTheme = Theme{
	Name: "Matrix Classic",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#003B00", ANSI256: "22", ANSI: "2"},
		Dark:  CompleteColor{TrueColor: "#003B00", ANSI256: "22", ANSI: "2"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#5FFF5F", ANSI256: "83", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#5FFF5F", ANSI256: "83", ANSI: "2"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00DD00", ANSI256: "40", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#00DD00", ANSI256: "40", ANSI: "2"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#5FFF5F", ANSI256: "83", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#5FFF5F", ANSI256: "83", ANSI: "2"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#003B00", ANSI256: "22", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#003B00", ANSI256: "22", ANSI: "2"},
		},
	},
}

// MatrixNeoTheme is a modern take on the matrix theme with blue accents
var MatrixNeoTheme = Theme{
	Name: "Matrix Neo",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
		Dark:  CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#5FFF5F", ANSI256: "83", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#5FFF5F", ANSI256: "83", ANSI: "2"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00EE00", ANSI256: "41", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#00EE00", ANSI256: "41", ANSI: "2"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#005555", ANSI256: "23", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#005555", ANSI256: "23", ANSI: "6"},
		},
	},
}

// CyberpunkNeonTheme is a vibrant cyberpunk theme with neon colors
var CyberpunkNeonTheme = Theme{
	Name: "Cyberpunk Neon",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
		Dark:  CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
		Dark:  CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#EE00EE", ANSI256: "200", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#EE00EE", ANSI256: "200", ANSI: "5"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#800080", ANSI256: "90", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#800080", ANSI256: "90", ANSI: "5"},
		},
	},
}

// CyberpunkRedTheme is a darker cyberpunk theme with red accents
var CyberpunkRedTheme = Theme{
	Name: "Cyberpunk Red",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FF0000", ANSI256: "196", ANSI: "1"},
		Dark:  CompleteColor{TrueColor: "#FF0000", ANSI256: "196", ANSI: "1"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
		Dark:  CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF0000", ANSI256: "196", ANSI: "1"},
			Dark:  CompleteColor{TrueColor: "#FF0000", ANSI256: "196", ANSI: "1"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF0000", ANSI256: "196", ANSI: "1"},
			Dark:  CompleteColor{TrueColor: "#FF0000", ANSI256: "196", ANSI: "1"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#EE0000", ANSI256: "196", ANSI: "1"},
			Dark:  CompleteColor{TrueColor: "#EE0000", ANSI256: "196", ANSI: "1"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF0000", ANSI256: "196", ANSI: "1"},
			Dark:  CompleteColor{TrueColor: "#FF0000", ANSI256: "196", ANSI: "1"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#8B0000", ANSI256: "88", ANSI: "1"},
			Dark:  CompleteColor{TrueColor: "#8B0000", ANSI256: "88", ANSI: "1"},
		},
	},
}

// PythonTheme uses Python's official colors
var PythonTheme = Theme{
	Name: "Python",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#306998", ANSI256: "25", ANSI: "4"},
		Dark:  CompleteColor{TrueColor: "#306998", ANSI256: "25", ANSI: "4"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FFE873", ANSI256: "221", ANSI: "3"},
		Dark:  CompleteColor{TrueColor: "#FFE873", ANSI256: "221", ANSI: "3"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#306998", ANSI256: "25", ANSI: "4"},
			Dark:  CompleteColor{TrueColor: "#306998", ANSI256: "25", ANSI: "4"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFE873", ANSI256: "221", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#FFE873", ANSI256: "221", ANSI: "3"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#306998", ANSI256: "25", ANSI: "4"},
			Dark:  CompleteColor{TrueColor: "#306998", ANSI256: "25", ANSI: "4"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#285E89", ANSI256: "24", ANSI: "4"},
			Dark:  CompleteColor{TrueColor: "#285E89", ANSI256: "24", ANSI: "4"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFE873", ANSI256: "221", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#FFE873", ANSI256: "221", ANSI: "3"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#306998", ANSI256: "25", ANSI: "4"},
			Dark:  CompleteColor{TrueColor: "#306998", ANSI256: "25", ANSI: "4"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFE873", ANSI256: "221", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#FFE873", ANSI256: "221", ANSI: "3"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#1B4B72", ANSI256: "24", ANSI: "4"},
			Dark:  CompleteColor{TrueColor: "#1B4B72", ANSI256: "24", ANSI: "4"},
		},
	},
}

// MonochromeTheme is a clean black and white theme
var MonochromeTheme = Theme{
	Name: "Monochrome",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
		Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		Dark:  CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
			Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
			Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
			Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
			Dark:  CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#666666", ANSI256: "241", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#999999", ANSI256: "247", ANSI: "7"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
			Dark:  CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#666666", ANSI256: "241", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#999999", ANSI256: "247", ANSI: "7"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#808080", ANSI256: "244", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#808080", ANSI256: "244", ANSI: "8"},
		},
	},
}

// RainbowBrightTheme is a vibrant rainbow theme
var RainbowBrightTheme = Theme{
	Name: "Rainbow Bright",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FF0000", ANSI256: "196", ANSI: "1"},
		Dark:  CompleteColor{TrueColor: "#FF0000", ANSI256: "196", ANSI: "1"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFFF00", ANSI256: "226", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#FFFF00", ANSI256: "226", ANSI: "3"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#EE00EE", ANSI256: "200", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#EE00EE", ANSI256: "200", ANSI: "5"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#4A4A4A", ANSI256: "239", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#4A4A4A", ANSI256: "239", ANSI: "8"},
		},
	},
}

// RainbowPastelTheme is a softer rainbow theme
var RainbowPastelTheme = Theme{
	Name: "Rainbow Pastel",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FFA3E7", ANSI256: "218", ANSI: "13"},
		Dark:  CompleteColor{TrueColor: "#FFA3E7", ANSI256: "218", ANSI: "13"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#BAFFC9", ANSI256: "157", ANSI: "10"},
		Dark:  CompleteColor{TrueColor: "#BAFFC9", ANSI256: "157", ANSI: "10"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#BAB3FF", ANSI256: "147", ANSI: "12"},
			Dark:  CompleteColor{TrueColor: "#BAB3FF", ANSI256: "147", ANSI: "12"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFA3F7", ANSI256: "219", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#FFA3F7", ANSI256: "219", ANSI: "13"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFFFBA", ANSI256: "229", ANSI: "11"},
			Dark:  CompleteColor{TrueColor: "#FFFFBA", ANSI256: "229", ANSI: "11"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFA3E7", ANSI256: "218", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#FFA3E7", ANSI256: "218", ANSI: "13"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#BAB3FF", ANSI256: "147", ANSI: "12"},
			Dark:  CompleteColor{TrueColor: "#BAB3FF", ANSI256: "147", ANSI: "12"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFA3F7", ANSI256: "219", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#FFA3F7", ANSI256: "219", ANSI: "13"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#BAFFC9", ANSI256: "157", ANSI: "10"},
			Dark:  CompleteColor{TrueColor: "#BAFFC9", ANSI256: "157", ANSI: "10"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#808080", ANSI256: "244", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#808080", ANSI256: "244", ANSI: "8"},
		},
	},
}

// BarbieTheme is inspired by the Barbie movie colors
var BarbieTheme = Theme{
	Name: "Barbie",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FF69B4", ANSI256: "205", ANSI: "13"},
		Dark:  CompleteColor{TrueColor: "#FF69B4", ANSI256: "205", ANSI: "13"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FFB6C1", ANSI256: "217", ANSI: "13"},
		Dark:  CompleteColor{TrueColor: "#FFB6C1", ANSI256: "217", ANSI: "13"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF1493", ANSI256: "198", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#FF1493", ANSI256: "198", ANSI: "5"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF69B4", ANSI256: "205", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#FF69B4", ANSI256: "205", ANSI: "13"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFB6C1", ANSI256: "217", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#FFB6C1", ANSI256: "217", ANSI: "13"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF59A4", ANSI256: "204", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#FF59A4", ANSI256: "204", ANSI: "13"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF1493", ANSI256: "198", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#FF1493", ANSI256: "198", ANSI: "5"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF69B4", ANSI256: "205", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#FF69B4", ANSI256: "205", ANSI: "13"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFB6C1", ANSI256: "217", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#FFB6C1", ANSI256: "217", ANSI: "13"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#C76E97", ANSI256: "168", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#C76E97", ANSI256: "168", ANSI: "5"},
		},
	},
}

// Retro Arcade Theme
var RetroArcadeTheme = Theme{
	Name: "Retro Arcade",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
		Dark:  CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
		Dark:  CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#EE00EE", ANSI256: "200", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#EE00EE", ANSI256: "200", ANSI: "5"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00FFFF", ANSI256: "51", ANSI: "6"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#800080", ANSI256: "90", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#800080", ANSI256: "90", ANSI: "5"},
		},
	},
}

// Forest Whisper Theme
var ForestWhisperTheme = Theme{
	Name: "Forest Whisper",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#228B22", ANSI256: "28", ANSI: "2"},
		Dark:  CompleteColor{TrueColor: "#228B22", ANSI256: "28", ANSI: "2"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
		Dark:  CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#A0522D", ANSI256: "130", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#A0522D", ANSI256: "130", ANSI: "3"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#228B22", ANSI256: "28", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#228B22", ANSI256: "28", ANSI: "2"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#1E7B1E", ANSI256: "28", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#1E7B1E", ANSI256: "28", ANSI: "2"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#A0522D", ANSI256: "130", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#A0522D", ANSI256: "130", ANSI: "3"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#228B22", ANSI256: "28", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#228B22", ANSI256: "28", ANSI: "2"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#2F4F2F", ANSI256: "22", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#2F4F2F", ANSI256: "22", ANSI: "2"},
		},
	},
}

// Ocean Breeze Theme
var OceanBreezeTheme = Theme{
	Name: "Ocean Breeze",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#00CED1", ANSI256: "44", ANSI: "6"},
		Dark:  CompleteColor{TrueColor: "#00CED1", ANSI256: "44", ANSI: "6"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00BFFF", ANSI256: "39", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00BFFF", ANSI256: "39", ANSI: "6"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00CED1", ANSI256: "44", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00CED1", ANSI256: "44", ANSI: "6"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
			Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00BEC1", ANSI256: "37", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00BEC1", ANSI256: "37", ANSI: "6"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00BFFF", ANSI256: "39", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00BFFF", ANSI256: "39", ANSI: "6"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00CED1", ANSI256: "44", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#00CED1", ANSI256: "44", ANSI: "6"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
			Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#006D6D", ANSI256: "23", ANSI: "6"},
			Dark:  CompleteColor{TrueColor: "#006D6D", ANSI256: "23", ANSI: "6"},
		},
	},
}

// Desert Sunset Theme
var DesertSunsetTheme = Theme{
	Name: "Desert Sunset",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FF4500", ANSI256: "202", ANSI: "9"},
		Dark:  CompleteColor{TrueColor: "#FF4500", ANSI256: "202", ANSI: "9"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
		Dark:  CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF6347", ANSI256: "196", ANSI: "9"},
			Dark:  CompleteColor{TrueColor: "#FF6347", ANSI256: "196", ANSI: "9"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF4500", ANSI256: "202", ANSI: "9"},
			Dark:  CompleteColor{TrueColor: "#FF4500", ANSI256: "202", ANSI: "9"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#EE3500", ANSI256: "202", ANSI: "9"},
			Dark:  CompleteColor{TrueColor: "#EE3500", ANSI256: "202", ANSI: "9"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF6347", ANSI256: "196", ANSI: "9"},
			Dark:  CompleteColor{TrueColor: "#FF6347", ANSI256: "196", ANSI: "9"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF4500", ANSI256: "202", ANSI: "9"},
			Dark:  CompleteColor{TrueColor: "#FF4500", ANSI256: "202", ANSI: "9"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
		},
	},
}

// Cyberpunk City Theme
var CyberpunkCityTheme = Theme{
	Name: "Cyberpunk City",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#8A2BE2", ANSI256: "93", ANSI: "5"},
		Dark:  CompleteColor{TrueColor: "#8A2BE2", ANSI256: "93", ANSI: "5"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#FF00FF", ANSI256: "201", ANSI: "5"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#8A2BE2", ANSI256: "93", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#8A2BE2", ANSI256: "93", ANSI: "5"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#7A1BD2", ANSI256: "92", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#7A1BD2", ANSI256: "92", ANSI: "5"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#8A2BE2", ANSI256: "93", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#8A2BE2", ANSI256: "93", ANSI: "5"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
			Dark:  CompleteColor{TrueColor: "#00FF00", ANSI256: "46", ANSI: "2"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#4B0082", ANSI256: "54", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#4B0082", ANSI256: "54", ANSI: "5"},
		},
	},
}

// Vintage Newspaper Theme
var VintageNewspaperTheme = Theme{
	Name: "Vintage Newspaper",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
		Dark:  CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#808080", ANSI256: "244", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#808080", ANSI256: "244", ANSI: "8"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
			Dark:  CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
			Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#1A1A1A", ANSI256: "234", ANSI: "0"},
			Dark:  CompleteColor{TrueColor: "#1A1A1A", ANSI256: "234", ANSI: "0"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#808080", ANSI256: "244", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#808080", ANSI256: "244", ANSI: "8"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
			Dark:  CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
			Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#696969", ANSI256: "242", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#696969", ANSI256: "242", ANSI: "8"},
		},
	},
}

// Steampunk Adventure Theme
var SteampunkAdventureTheme = Theme{
	Name: "Steampunk Adventure",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#A8760B", ANSI256: "136", ANSI: "3"},
		Dark:  CompleteColor{TrueColor: "#A8760B", ANSI256: "136", ANSI: "3"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
		Dark:  CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#A8760B", ANSI256: "136", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#A8760B", ANSI256: "136", ANSI: "3"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#A8760B", ANSI256: "136", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#A8760B", ANSI256: "136", ANSI: "3"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#A8760B", ANSI256: "136", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#A8760B", ANSI256: "136", ANSI: "3"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#8B4513", ANSI256: "94", ANSI: "3"},
		},
	},
}

// Galaxy Quest Theme
var GalaxyQuestTheme = Theme{
	Name: "Galaxy Quest",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#3B0072", ANSI256: "54", ANSI: "5"},
		Dark:  CompleteColor{TrueColor: "#3B0072", ANSI256: "54", ANSI: "5"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
		Dark:  CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#3B0072", ANSI256: "54", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#3B0072", ANSI256: "54", ANSI: "5"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#3B0072", ANSI256: "54", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#3B0072", ANSI256: "54", ANSI: "5"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
			Dark:  CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#3B0072", ANSI256: "54", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#3B0072", ANSI256: "54", ANSI: "5"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
			Dark:  CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#3B0072", ANSI256: "54", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#3B0072", ANSI256: "54", ANSI: "5"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
			Dark:  CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#2E1A47", ANSI256: "54", ANSI: "5"},
			Dark:  CompleteColor{TrueColor: "#2E1A47", ANSI256: "54", ANSI: "5"},
		},
	},
}

// Minimalist Zen Theme
var MinimalistZenTheme = Theme{
	Name: "Minimalist Zen",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#707070", ANSI256: "243", ANSI: "8"},
		Dark:  CompleteColor{TrueColor: "#707070", ANSI256: "243", ANSI: "8"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#707070", ANSI256: "243", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#707070", ANSI256: "243", ANSI: "8"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#707070", ANSI256: "243", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#707070", ANSI256: "243", ANSI: "8"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
			Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#707070", ANSI256: "243", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#707070", ANSI256: "243", ANSI: "8"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
			Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#707070", ANSI256: "243", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#707070", ANSI256: "243", ANSI: "8"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
			Dark:  CompleteColor{TrueColor: "#FFFFFF", ANSI256: "231", ANSI: "15"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#A9A9A9", ANSI256: "248", ANSI: "8"},
			Dark:  CompleteColor{TrueColor: "#A9A9A9", ANSI256: "248", ANSI: "8"},
		},
	},
}

// Candy Land Theme
var CandyLandTheme = Theme{
	Name: "Candy Land",
	Primary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FF59A4", ANSI256: "204", ANSI: "13"},
		Dark:  CompleteColor{TrueColor: "#FF59A4", ANSI256: "204", ANSI: "13"},
	},
	Secondary: CompleteAdaptiveColor{
		Light: CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
		Dark:  CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
	},
	StatusBar: StatusBarColors{
		Text: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF59A4", ANSI256: "204", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#FF59A4", ANSI256: "204", ANSI: "13"},
		},
		Title: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
		},
		Model: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF59A4", ANSI256: "204", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#FF59A4", ANSI256: "204", ANSI: "13"},
		},
	},
	Border: BorderColors{
		Normal: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF59A4", ANSI256: "204", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#FF59A4", ANSI256: "204", ANSI: "13"},
		},
		Active: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
		},
	},
	Message: MessageColors{
		UserText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FF59A4", ANSI256: "204", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#FF59A4", ANSI256: "204", ANSI: "13"},
		},
		AIText: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
			Dark:  CompleteColor{TrueColor: "#FFD700", ANSI256: "220", ANSI: "3"},
		},
		Timestamp: CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: "#DDA0DD", ANSI256: "182", ANSI: "13"},
			Dark:  CompleteColor{TrueColor: "#DDA0DD", ANSI256: "182", ANSI: "13"},
		},
	},
}

func LoadThemeFromConfig(themeName string) {
	switch themeName {
	case "Default":
		CurrentTheme = DefaultTheme
	case "Dracula":
		CurrentTheme = DraculaTheme
	case "Nord":
		CurrentTheme = NordTheme
	case "Matrix Classic":
		CurrentTheme = MatrixClassicTheme
	case "Matrix Neo":
		CurrentTheme = MatrixNeoTheme
	case "Cyberpunk Neon":
		CurrentTheme = CyberpunkNeonTheme
	case "Cyberpunk Red":
		CurrentTheme = CyberpunkRedTheme
	case "Python":
		CurrentTheme = PythonTheme
	case "Monochrome":
		CurrentTheme = MonochromeTheme
	case "Rainbow Bright":
		CurrentTheme = RainbowBrightTheme
	case "Rainbow Pastel":
		CurrentTheme = RainbowPastelTheme
	case "Barbie":
		CurrentTheme = BarbieTheme
	default:
		CurrentTheme = DefaultTheme
	}
}