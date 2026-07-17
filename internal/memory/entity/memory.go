package entity

import "time"

type Memory struct {
	ID         string    `gorm:"type:char(26);primaryKey"`
	VisitorID  string    `gorm:"type:char(26);not null;index"`
	Category   string    `gorm:"type:text"`
	MemoryText string    `gorm:"type:text;not null"`
	Importance int       `gorm:"type:int;default:3"`
	CreatedAt  time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt  time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

func (Memory) TableName() string {
	return "visitor_knowledge"
}
