package entity

import "time"

type Memory struct {
	ID              string     `gorm:"type:char(26);primaryKey"`
	VisitorID       string     `gorm:"type:char(26);not null;index:idx_visitor_key,unique"`
	Category        string     `gorm:"type:text"`
	Key             string     `gorm:"type:text;not null;index:idx_visitor_key,unique"`
	Value           string     `gorm:"type:text;not null"`
	Confidence      float32    `gorm:"type:numeric(5,4);default:0.75"`
	LastConfirmedAt *time.Time `gorm:"type:timestamptz"`
	CreatedAt       time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt       time.Time  `gorm:"type:timestamptz;not null;default:now()"`
}

func (Memory) TableName() string {
	return "visitor_memories"
}
