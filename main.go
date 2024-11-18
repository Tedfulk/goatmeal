package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"goatmeal/config"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	list     list.Model
	quitting bool
}

func main() {
	// Check for the configuration file
	usr, _ := user.Current()
	configPath := filepath.Join(usr.HomeDir, ".goatmeal", "config.yaml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Run the wizard if the config file does not exist
		runWizard()
	} else {
		fmt.Println("Configuration file found. Proceeding with the application...")
		// Proceed with the rest of the application
	}
}

func runWizard() {
	// Start the Bubble Tea application in full-screen mode
	p := tea.NewProgram(config.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
