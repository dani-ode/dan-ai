package publisher

import (
	"context"
	"log"
	"portfolio-ai/internal/outbox/repository"
	"portfolio-ai/pkg/kafka"
	"time"
)

const TopicPortfolioKnowledge = "portfolio.knowledge"
const TopicPortfolioEvents = "portfolio.events"

type Publisher struct {
	repo          repository.Repository
	kafkaProducer kafka.Producer
}

func NewPublisher(repo repository.Repository, kafkaProducer kafka.Producer) *Publisher {
	return &Publisher{
		repo:          repo,
		kafkaProducer: kafkaProducer,
	}
}

func (p *Publisher) Start(ctx context.Context, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	log.Println("Outbox publisher started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Outbox publisher stopped")
			return
		case <-ticker.C:
			p.processOutbox(ctx)
		}
	}
}

func (p *Publisher) processOutbox(ctx context.Context) {
	events, err := p.repo.GetUnpublishedEvents(ctx, 50)
	if err != nil {
		log.Printf("error fetching unpublished events: %v", err)
		return
	}

	for _, evt := range events {
		// Map OutboxEvent to Kafka Event
		kEvent := kafka.Event{
			Aggregate:   evt.Aggregate,
			AggregateID: evt.AggregateID,
			EventType:   evt.EventType,
			Payload:     evt.Payload,
			Timestamp:   evt.CreatedAt,
		}

		// Determine target topic
		var targetTopic string
		if evt.Aggregate == "knowledge_document" || evt.Aggregate == "project" || evt.Aggregate == "profile" || evt.Aggregate == "experience" || evt.Aggregate == "certificate" {
			targetTopic = TopicPortfolioKnowledge
		} else {
			targetTopic = TopicPortfolioEvents
		}

		// Publish to Kafka
		if err := p.kafkaProducer.Publish(ctx, targetTopic, kEvent); err != nil {
			log.Printf("failed to publish event %s: %v", evt.ID, err)
			_ = p.repo.MarkAsFailed(ctx, evt.ID, err.Error())
			continue
		}

		// Mark as published
		if err := p.repo.MarkAsPublished(ctx, evt.ID); err != nil {
			log.Printf("failed to mark event %s as published: %v", evt.ID, err)
		}
	}
}
