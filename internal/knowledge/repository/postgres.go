package repository

import (
	"context"

	"portfolio-ai/internal/knowledge/entity"
	outboxEntity "portfolio-ai/internal/outbox/entity"
	"gorm.io/gorm"
)

type KnowledgeRepository interface {
	// Document
	CreateDocument(ctx context.Context, doc *entity.KnowledgeDocument) error
	UpdateDocument(ctx context.Context, doc *entity.KnowledgeDocument) error
	DeleteDocument(ctx context.Context, id string) error
	GetDocumentByID(ctx context.Context, id string) (*entity.KnowledgeDocument, error)
	GetDocumentBySource(ctx context.Context, sourceType, sourceID string) (*entity.KnowledgeDocument, error)
	ListDocuments(ctx context.Context, page, pageSize int, sourceType string) ([]entity.KnowledgeDocument, int64, error)

	// Chunk
	CreateChunks(ctx context.Context, chunks []entity.KnowledgeChunk) error
	DeleteChunksByDocumentID(ctx context.Context, documentID string) error
	ListChunksByDocumentID(ctx context.Context, documentID string) ([]entity.KnowledgeChunk, error)

	// Outbox
	CreateOutboxEvent(ctx context.Context, event *outboxEntity.OutboxEvent) error
}

type postgresKnowledgeRepository struct {
	db *gorm.DB
}

func NewPostgresKnowledgeRepository(db *gorm.DB) KnowledgeRepository {
	return &postgresKnowledgeRepository{db: db}
}

// Document

func (r *postgresKnowledgeRepository) CreateDocument(ctx context.Context, doc *entity.KnowledgeDocument) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *postgresKnowledgeRepository) UpdateDocument(ctx context.Context, doc *entity.KnowledgeDocument) error {
	return r.db.WithContext(ctx).Save(doc).Error
}

func (r *postgresKnowledgeRepository) DeleteDocument(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.KnowledgeDocument{}, "id = ?", id).Error
}

func (r *postgresKnowledgeRepository) GetDocumentByID(ctx context.Context, id string) (*entity.KnowledgeDocument, error) {
	var doc entity.KnowledgeDocument
	err := r.db.WithContext(ctx).First(&doc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *postgresKnowledgeRepository) GetDocumentBySource(ctx context.Context, sourceType, sourceID string) (*entity.KnowledgeDocument, error) {
	var doc entity.KnowledgeDocument
	err := r.db.WithContext(ctx).First(&doc, "source_type = ? AND source_id = ?", sourceType, sourceID).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *postgresKnowledgeRepository) ListDocuments(ctx context.Context, page, pageSize int, sourceType string) ([]entity.KnowledgeDocument, int64, error) {
	var docs []entity.KnowledgeDocument
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.KnowledgeDocument{})
	if sourceType != "" {
		query = query.Where("source_type = ?", sourceType)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&docs).Error
	return docs, total, err
}

// Chunk

func (r *postgresKnowledgeRepository) CreateChunks(ctx context.Context, chunks []entity.KnowledgeChunk) error {
	if len(chunks) == 0 {
		return nil
	}
	// batch insert
	return r.db.WithContext(ctx).Create(&chunks).Error
}

func (r *postgresKnowledgeRepository) DeleteChunksByDocumentID(ctx context.Context, documentID string) error {
	return r.db.WithContext(ctx).Delete(&entity.KnowledgeChunk{}, "document_id = ?", documentID).Error
}

func (r *postgresKnowledgeRepository) ListChunksByDocumentID(ctx context.Context, documentID string) ([]entity.KnowledgeChunk, error) {
	var chunks []entity.KnowledgeChunk
	err := r.db.WithContext(ctx).Where("document_id = ?", documentID).Order("chunk_index ASC").Find(&chunks).Error
	return chunks, err
}

// Outbox

func (r *postgresKnowledgeRepository) CreateOutboxEvent(ctx context.Context, event *outboxEntity.OutboxEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}
