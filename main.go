package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/database"
	"github.com/tedfulk/goatmeal/ui"
	"github.com/tedfulk/goatmeal/ui/setup"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Check if setup wizard needs to run
	if cfg.CurrentModel == "" {
		wizard := setup.NewWizard(cfg)
		p := tea.NewProgram(wizard)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running setup wizard: %v\n", err)
			os.Exit(1)
		}

		// Reload config after wizard completes
		cfg, err = config.Load()
		if err != nil {
			fmt.Printf("Error loading config after setup: %v\n", err)
			os.Exit(1)
		}

		// Check if setup was completed
		if cfg.CurrentModel == "" {
			fmt.Println("Setup wizard was not completed. Please run the program again to complete setup.")
			os.Exit(1)
		}
	}

	// Initialize database
	dbPath := os.ExpandEnv("$HOME/.config/goatmeal/goatmeal.db")
	db, err := database.NewDB(dbPath)
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Clean up old conversations
	if err := db.CleanupOldConversations(cfg.Settings.ConversationRetention); err != nil {
		fmt.Printf("Error cleaning up old conversations: %v\n", err)
	}

	// Initialize UI
	app := ui.NewApp(cfg)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		os.Exit(1)
	}
} 