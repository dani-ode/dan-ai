package entity

import (
	"time"
)

type KnowledgeDocument struct {
	ID             string         `gorm:"primaryKey;type:varchar(26)" json:"id"` // ULID
	SourceType     string         `gorm:"type:varchar(50);not null" json:"source_type"`
	SourceID       string         `gorm:"type:varchar(26);not null" json:"source_id"`
	Title          string         `gorm:"type:varchar(255);not null" json:"title"`
	Content        string         `gorm:"type:text;not null" json:"content"`
	Checksum       string         `gorm:"type:varchar(64);not null" json:"checksum"`
	Version        int32          `gorm:"not null;default:1" json:"version"`
	Status         string         `gorm:"type:varchar(20);not null;default:'Pending'" json:"status"`
	EmbeddingModel string         `gorm:"type:varchar(50)" json:"embedding_model"`
	LastEmbeddedAt *time.Time     `json:"last_embedded_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
