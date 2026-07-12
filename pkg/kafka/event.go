package kafka

import (
	"encoding/json"
	"time"
)

// Event represents the schema of messages sent to Kafka
type Event struct {
	Aggregate   string          `json:"aggregate"`
	AggregateID string          `json:"aggregate_id"`
	EventType   string          `json:"event_type"`
	Payload     json.RawMessage `json:"payload"`
	Timestamp   time.Time       `json:"timestamp"`
}
