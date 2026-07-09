// internal/chat/entity/session.go
package entity

import "time"

type ChatSession struct {
	ID        string     `gorm:"type:char(26);primaryKey"`
	VisitorID string     `gorm:"type:char(26);not null"`
	PromptID  string     `gorm:"type:char(26)"`
	Title     string     `gorm:"type:text"`
	StartedAt time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	EndedAt   *time.Time `gorm:"type:timestamptz"`
	CreatedAt time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time  `gorm:"type:timestamptz;not null;default:now()"`
}

// TableName overrides the table name used by ChatSession to chat_sessions.
func (ChatSession) TableName() string {
	return "chat_sessions"
}
