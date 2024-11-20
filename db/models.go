package db

import (
	"time"
)

type Conversation struct {
    ID        string
    Title     string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Message struct {
    ID             string
    ConversationID string
    Role           string
    Content        string
    CreatedAt      time.Time
}

type ChatDB interface {
    // Conversation methods
    CreateConversation() (string, error)
    GetConversation(id string) (*Conversation, error)
    ListConversations() ([]Conversation, error)
    UpdateConversationTitle(id string, title string) error
    
    // Message methods
    AddMessage(convID string, msg Message) error
    GetMessages(convID string) ([]Message, error)
} 