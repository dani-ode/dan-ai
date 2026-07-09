// internal/chat/entity/message.go
package entity

import "time"

type ChatMessage struct {
	ID               string    `gorm:"type:char(26);primaryKey"`
	SessionID        string    `gorm:"type:char(26);not null"`
	Role             string    `gorm:"type:text;not null"`
	Content          string    `gorm:"type:text;not null"`
	Model            string    `gorm:"type:text"`
	PromptTokens     int32     `gorm:"type:int"`
	CompletionTokens int32     `gorm:"type:int"`
	LatencyMs        int32     `gorm:"type:int"`
	Status           string    `gorm:"type:text;not null;default:'Pending'"`
	CreatedAt        time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt        time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

// TableName overrides the table name used by ChatMessage to chat_messages.
func (ChatMessage) TableName() string {
	return "chat_messages"
}
