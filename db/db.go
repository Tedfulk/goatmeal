package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ChatDB struct {
	db *sql.DB
}

func NewChatDB() (*ChatDB, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %w", err)
	}

	// Create .goatmeal/db directory if it doesn't exist
	dbDir := filepath.Join(homeDir, ".goatmeal", "db")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating database directory: %w", err)
	}

	// Set up database path
	dbPath := filepath.Join(dbDir, "chat.db")

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Create tables if they don't exist
	if err := initializeTables(db); err != nil {
		db.Close()
		return nil, err
	}

	// Set up automatic cleanup trigger
	if err := setupCleanupTrigger(db); err != nil {
		db.Close()
		return nil, err
	}

	return &ChatDB{db: db}, nil
}

func initializeTables(db *sql.DB) error {
	// Create messages table with timestamp
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			role TEXT NOT NULL,
			content TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func setupCleanupTrigger(db *sql.DB) error {
	// Create trigger to automatically delete old messages
	_, err := db.Exec(`
		CREATE TRIGGER IF NOT EXISTS cleanup_old_messages
		AFTER INSERT ON messages
		BEGIN
			DELETE FROM messages 
			WHERE timestamp < datetime('now', '-30 days');
		END;
	`)
	return err
}

func (c *ChatDB) SaveMessage(role, content string) error {
	_, err := c.db.Exec(`
		INSERT INTO messages (role, content, timestamp)
		VALUES (?, ?, datetime('now'))
	`, role, content)
	return err
}

func (c *ChatDB) GetRecentMessages(limit int) ([]Message, error) {
	rows, err := c.db.Query(`
		SELECT role, content, timestamp
		FROM messages
		ORDER BY timestamp DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var timestamp time.Time
		if err := rows.Scan(&msg.Role, &msg.Content, &timestamp); err != nil {
			return nil, err
		}
		msg.Timestamp = timestamp
		messages = append(messages, msg)
	}
	return messages, nil
}

type Message struct {
	Role      string
	Content   string
	Timestamp time.Time
}

func (c *ChatDB) Close() error {
	return c.db.Close()
} 