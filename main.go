package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"goatmeal/config"
	"goatmeal/db"
	"goatmeal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Check for the configuration file
	usr, _ := user.Current()
	configPath := filepath.Join(usr.HomeDir, ".goatmeal", "config.yaml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Run the wizard if the config file does not exist
		if err := runWizard(); err != nil {
			fmt.Printf("Error during setup: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize database
	chatDB, err := db.NewChatDB()
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer chatDB.Close()

	// Start the chat application with the config and database
	if err := runChat(chatDB); err != nil {
		fmt.Printf("Error running chat: %v\n", err)
		os.Exit(1)
	}
}

func runWizard() error {
	p := tea.NewProgram(config.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running wizard: %w", err)
	}
	return nil
}

func runChat(chatDB db.ChatDB) error {
	model, err := ui.NewMainModel(chatDB)
	if err != nil {
		return fmt.Errorf("error initializing chat: %w", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running chat: %w", err)
	}

	return nil
}
