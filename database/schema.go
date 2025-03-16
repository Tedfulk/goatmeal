package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS conversations (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    provider TEXT NOT NULL,
    model TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_conversations_created_at ON conversations(created_at);
`

// DB represents the database connection
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection
func NewDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("error creating schema: %w", err)
	}

	return &DB{db}, nil
}

// CleanupOldConversations deletes conversations older than the retention period
func (db *DB) CleanupOldConversations(retentionDays int) error {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	
	_, err := db.Exec(`
		DELETE FROM conversations 
		WHERE created_at < ?
	`, cutoff)
	
	if err != nil {
		return fmt.Errorf("error cleaning up old conversations: %w", err)
	}
	
	return nil
}

// Message represents a chat message
type Message struct {
	ID             string
	ConversationID string
	Role           string    // Can be "user", "assistant", or "search"
	Content        string
	CreatedAt      time.Time
}

// Conversation represents a chat conversation
type Conversation struct {
	ID        string
	Title     string
	Provider  string
	Model     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Messages  []Message
} 