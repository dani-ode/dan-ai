package entity

import (
	"encoding/json"
	"time"
)

type OutboxEvent struct {
	ID           string          `gorm:"primaryKey;type:varchar(26)" json:"id"` // ULID
	Aggregate    string          `gorm:"type:text;not null" json:"aggregate"`
	AggregateID  string          `gorm:"type:varchar(26);not null" json:"aggregate_id"`
	EventType    string          `gorm:"type:text;not null" json:"event_type"`
	Payload      json.RawMessage `gorm:"type:jsonb;not null" json:"payload"`
	Published    bool            `gorm:"default:false" json:"published"`
	RetryCount   int             `gorm:"default:0" json:"retry_count"`
	FailedReason *string         `gorm:"type:text" json:"failed_reason"`
	PublishedAt  *time.Time      `json:"published_at"`
	CreatedAt    time.Time       `json:"created_at"`
}
