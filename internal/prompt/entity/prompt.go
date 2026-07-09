// internal/prompt/entity/prompt.go
package entity

import "time"

type Prompt struct {
	ID           string    `gorm:"type:char(26);primaryKey"`
	Name         string    `gorm:"type:text;uniqueIndex;not null"`
	SystemPrompt string    `gorm:"type:text;not null"`
	Description  string    `gorm:"type:text"`
	ModelID      string    `gorm:"type:char(26)"`
	Active       bool      `gorm:"type:boolean;default:false"`
	Version      int32     `gorm:"type:int;default:1"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

// TableName overrides the table name used by Prompt to prompts.
func (Prompt) TableName() string {
	return "prompts"
}
