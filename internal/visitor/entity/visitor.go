// internal/visitor/entity/visitor.go
package entity

import "time"

type Visitor struct {
	ID            string    `gorm:"type:char(26);primaryKey"`
	FirstSeenAt   time.Time `gorm:"type:timestamptz;not null;default:now()"`
	LastSeenAt    time.Time `gorm:"type:timestamptz;not null;default:now()"`
	TotalMessages int32     `gorm:"type:int;default:0"`
	CreatedAt     time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt     time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

// TableName overrides the table name used by Visitor to visitors.
func (Visitor) TableName() string {
	return "visitors"
}
