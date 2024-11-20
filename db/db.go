package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteChatDB struct {
	db *sql.DB
}

func NewChatDB() (*SQLiteChatDB, error) {
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

	// Create tables
	if err := initializeTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return &SQLiteChatDB{db: db}, nil
}

func initializeTables(db *sql.DB) error {
	// Create conversations table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS conversations (
			id TEXT PRIMARY KEY,
			title TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating conversations table: %w", err)
	}

	// Create messages table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			conversation_id TEXT NOT NULL,
			role TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			FOREIGN KEY (conversation_id) REFERENCES conversations(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating messages table: %w", err)
	}

	return nil
}

// Implement ChatDB interface methods

func (db *SQLiteChatDB) CreateConversation() (string, error) {
	id := uuid.New().String()
	now := time.Now()

	_, err := db.db.Exec(`
		INSERT INTO conversations (id, created_at, updated_at)
		VALUES (?, ?, ?)
	`, id, now, now)

	if err != nil {
		return "", fmt.Errorf("error creating conversation: %w", err)
	}

	return id, nil
}

func (db *SQLiteChatDB) GetConversation(id string) (*Conversation, error) {
	conv := &Conversation{}
	err := db.db.QueryRow(`
		SELECT id, title, created_at, updated_at
		FROM conversations
		WHERE id = ?
	`, id).Scan(&conv.ID, &conv.Title, &conv.CreatedAt, &conv.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting conversation: %w", err)
	}

	return conv, nil
}

func (db *SQLiteChatDB) ListConversations() ([]Conversation, error) {
	rows, err := db.db.Query(`
		SELECT id, title, created_at, updated_at
		FROM conversations
		ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("error listing conversations: %w", err)
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		var title sql.NullString
		err := rows.Scan(&conv.ID, &title, &conv.CreatedAt, &conv.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning conversation: %w", err)
		}
		
		if !title.Valid || title.String == "" {
			conv.Title = fmt.Sprintf("Conversation %d", len(conversations)+1)
		} else {
			conv.Title = title.String
		}
		
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

func (db *SQLiteChatDB) AddMessage(convID string, msg Message) error {
	msg.ID = uuid.New().String()
	msg.CreatedAt = time.Now()

	_, err := db.db.Exec(`
		INSERT INTO messages (id, conversation_id, role, content, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, msg.ID, convID, msg.Role, msg.Content, msg.CreatedAt)

	if err != nil {
		return fmt.Errorf("error adding message: %w", err)
	}

	// Update conversation's updated_at timestamp
	_, err = db.db.Exec(`
		UPDATE conversations
		SET updated_at = ?
		WHERE id = ?
	`, msg.CreatedAt, convID)

	if err != nil {
		return fmt.Errorf("error updating conversation timestamp: %w", err)
	}

	return nil
}

func (db *SQLiteChatDB) GetMessages(convID string) ([]Message, error) {
	rows, err := db.db.Query(`
		SELECT id, conversation_id, role, content, created_at
		FROM messages
		WHERE conversation_id = ?
		ORDER BY created_at ASC
	`, convID)
	if err != nil {
		return nil, fmt.Errorf("error getting messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (db *SQLiteChatDB) Close() error {
	return db.db.Close()
}

func (db *SQLiteChatDB) UpdateConversationTitle(id string, title string) error {
	_, err := db.db.Exec(`
		UPDATE conversations
		SET title = ?
		WHERE id = ?
	`, title, id)

	if err != nil {
		return fmt.Errorf("error updating conversation title: %w", err)
	}

	return nil
} 