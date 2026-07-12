package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"portfolio-ai/pkg/config"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Publish(ctx context.Context, topic string, event Event) error
	Close() error
}

type producer struct {
	writers map[string]*kafka.Writer
	brokers []string
}

func NewProducer(cfg *config.Config) Producer {
	return &producer{
		writers: make(map[string]*kafka.Writer),
		brokers: cfg.Kafka.Brokers,
	}
}

func (p *producer) getWriter(topic string) *kafka.Writer {
	if w, ok := p.writers[topic]; ok {
		return w
	}

	w := &kafka.Writer{
		Addr:                   kafka.TCP(p.brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	p.writers[topic] = w
	return w
}

func (p *producer) Publish(ctx context.Context, topic string, event Event) error {
	w := p.getWriter(topic)

	b, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(event.AggregateID),
		Value: b,
		Time:  event.Timestamp,
	}

	if err := w.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

func (p *producer) Close() error {
	for _, w := range p.writers {
		_ = w.Close()
	}
	return nil
}
