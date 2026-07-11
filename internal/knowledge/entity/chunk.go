package entity

import (
	"time"
)

type KnowledgeChunk struct {
	ID             string    `gorm:"primaryKey;type:varchar(26)" json:"id"` // ULID
	DocumentID     string    `gorm:"type:varchar(26);not null;index" json:"document_id"`
	ChunkIndex     int32     `gorm:"not null" json:"chunk_index"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	TokenCount     int32     `gorm:"not null" json:"token_count"`
	EmbeddingModel string    `gorm:"type:varchar(50)" json:"embedding_model"`
	CreatedAt      time.Time `json:"created_at"`
}
