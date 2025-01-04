package database

import (
	"fmt"
	"time"
)

// GetConversations retrieves a paginated list of conversations
// If limit is -1, returns all conversations
func (db *DB) GetConversations(offset, limit int) ([]Conversation, error) {
	var query string
	var args []interface{}

	if limit == -1 {
		query = `
			SELECT id, title, provider, model, created_at, updated_at
			FROM conversations
			ORDER BY updated_at DESC
		`
	} else {
		query = `
			SELECT id, title, provider, model, created_at, updated_at
			FROM conversations
			ORDER BY updated_at DESC
			LIMIT ? OFFSET ?
		`
		args = []interface{}{limit, offset}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying conversations: %w", err)
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(
			&conv.ID,
			&conv.Title,
			&conv.Provider,
			&conv.Model,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning conversation: %w", err)
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

// GetConversationMessages retrieves all messages for a conversation
func (db *DB) GetConversationMessages(conversationID string) ([]Message, error) {
	rows, err := db.Query(`
		SELECT id, role, content, created_at
		FROM messages
		WHERE conversation_id = ?
		ORDER BY created_at ASC
	`, conversationID)
	if err != nil {
		return nil, fmt.Errorf("error querying messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(
			&msg.ID,
			&msg.Role,
			&msg.Content,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning message: %w", err)
		}
		msg.ConversationID = conversationID
		messages = append(messages, msg)
	}

	return messages, nil
}

// SaveConversation saves a new conversation and its messages
func (db *DB) SaveConversation(conv *Conversation) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if conversation exists
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM conversations WHERE id = ?)", conv.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking conversation existence: %w", err)
	}

	if !exists {
		// Insert new conversation
		_, err = tx.Exec(`
			INSERT INTO conversations (id, title, provider, model, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, conv.ID, conv.Title, conv.Provider, conv.Model, conv.CreatedAt, conv.UpdatedAt)
		if err != nil {
			return fmt.Errorf("error inserting conversation: %w", err)
		}
	}

	// Update conversation's updated_at timestamp
	_, err = tx.Exec(`
		UPDATE conversations
		SET updated_at = ?
		WHERE id = ?
	`, time.Now(), conv.ID)
	if err != nil {
		return fmt.Errorf("error updating conversation timestamp: %w", err)
	}

	// Insert messages
	for _, msg := range conv.Messages {
		_, err = tx.Exec(`
			INSERT INTO messages (id, conversation_id, role, content, created_at)
			VALUES (?, ?, ?, ?, ?)
		`, msg.ID, conv.ID, msg.Role, msg.Content, msg.CreatedAt)
		if err != nil {
			return fmt.Errorf("error inserting message: %w", err)
		}
	}

	return tx.Commit()
}

// AddMessage adds a single message to an existing conversation
func (db *DB) AddMessage(msg *Message) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	// Update conversation's updated_at timestamp
	_, err = tx.Exec(`
		UPDATE conversations
		SET updated_at = ?
		WHERE id = ?
	`, time.Now(), msg.ConversationID)
	if err != nil {
		return fmt.Errorf("error updating conversation timestamp: %w", err)
	}

	// Insert the message
	_, err = tx.Exec(`
		INSERT INTO messages (id, conversation_id, role, content, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, msg.ID, msg.ConversationID, msg.Role, msg.Content, msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("error inserting message: %w", err)
	}

	return tx.Commit()
}

// UpdateConversationTitle updates the title of a conversation
func (db *DB) UpdateConversationTitle(conversationID, title string) error {
	_, err := db.Exec(`
		UPDATE conversations
		SET title = ?, updated_at = ?
		WHERE id = ?
	`, title, time.Now(), conversationID)
	
	if err != nil {
		return fmt.Errorf("error updating conversation title: %w", err)
	}
	
	return nil
}

// DeleteConversation deletes a conversation and its messages
func (db *DB) DeleteConversation(conversationID string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete messages first due to foreign key constraint
	_, err = tx.Exec(`DELETE FROM messages WHERE conversation_id = ?`, conversationID)
	if err != nil {
		return fmt.Errorf("error deleting messages: %w", err)
	}

	// Delete conversation
	_, err = tx.Exec(`DELETE FROM conversations WHERE id = ?`, conversationID)
	if err != nil {
		return fmt.Errorf("error deleting conversation: %w", err)
	}

	return tx.Commit()
} 