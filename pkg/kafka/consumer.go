package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"portfolio-ai/pkg/config"

	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	Consume(ctx context.Context, handler func(ctx context.Context, event Event) error) error
	Close() error
}

type consumer struct {
	reader *kafka.Reader
}

func NewConsumer(cfg *config.Config, topic string, groupID string) Consumer {
	// Ensure topic exists to avoid hanging reader
	if len(cfg.Kafka.Brokers) > 0 {
		for {
			conn, err := kafka.Dial("tcp", cfg.Kafka.Brokers[0])
			if err == nil {
				err = conn.CreateTopics(kafka.TopicConfig{
					Topic:             topic,
					NumPartitions:     1,
					ReplicationFactor: 1,
				})
				_ = conn.Close()
				if err == nil {
					log.Printf("topic %s created or already exists", topic)
					break
				}
				log.Printf("failed to create topic: %v, retrying...", err)
			} else {
				log.Printf("failed to dial kafka for topic creation: %v, retrying...", err)
			}
			time.Sleep(2 * time.Second)
		}
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Kafka.Brokers,
		GroupID:     groupID,
		Topic:       topic,
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.FirstOffset,
	})

	return &consumer{
		reader: r,
	}
}

func (c *consumer) Consume(ctx context.Context, handler func(ctx context.Context, event Event) error) error {
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Printf("error fetching message: %v", err)
			continue
		}

		var event Event
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("error unmarshaling event: %v", err)
			// Commit the message to skip it and avoid infinite loop
			_ = c.reader.CommitMessages(ctx, m)
			continue
		}

		if err := handler(ctx, event); err != nil {
			log.Printf("error handling event %s: %v", event.AggregateID, err)
			// Depending on retry policy, we could choose not to commit
			// But for now, let's commit so it moves to the next message
		}

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Printf("error committing message: %v", err)
		}
	}
}

func (c *consumer) Close() error {
	return c.reader.Close()
}
