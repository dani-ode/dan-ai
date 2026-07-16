// internal/chat/mapper/mapper.go
package mapper

import (
	"dan-ai/internal/chat/entity"
	pb "dan-ai/proto/chat"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// --- Session mappers ---

// SessionToProto maps a domain ChatSession entity to a protobuf ChatSession message.
func SessionToProto(e *entity.ChatSession) *pb.ChatSession {
	if e == nil {
		return nil
	}
	session := &pb.ChatSession{
		Id:        e.ID,
		VisitorId: e.VisitorID,
		PromptId:  e.PromptID,
		Title:     e.Title,
		StartedAt: timestamppb.New(e.StartedAt),
		CreatedAt: timestamppb.New(e.CreatedAt),
		UpdatedAt: timestamppb.New(e.UpdatedAt),
	}
	if e.EndedAt != nil {
		session.EndedAt = timestamppb.New(*e.EndedAt)
	}
	return session
}

// SessionsToProto maps a slice of ChatSession entities to protobuf messages.
func SessionsToProto(entities []entity.ChatSession) []*pb.ChatSession {
	protos := make([]*pb.ChatSession, len(entities))
	for i := range entities {
		protos[i] = SessionToProto(&entities[i])
	}
	return protos
}

// --- Message mappers ---

// MessageToProto maps a domain ChatMessage entity to a protobuf ChatMessage message.
func MessageToProto(e *entity.ChatMessage) *pb.ChatMessage {
	if e == nil {
		return nil
	}
	return &pb.ChatMessage{
		Id:               e.ID,
		SessionId:        e.SessionID,
		Role:             e.Role,
		Content:          e.Content,
		Model:            e.Model,
		PromptTokens:     e.PromptTokens,
		CompletionTokens: e.CompletionTokens,
		LatencyMs:        e.LatencyMs,
		Status:           e.Status,
		CreatedAt:        timestamppb.New(e.CreatedAt),
		UpdatedAt:        timestamppb.New(e.UpdatedAt),
	}
}

// MessagesToProto maps a slice of ChatMessage entities to protobuf messages.
func MessagesToProto(entities []entity.ChatMessage) []*pb.ChatMessage {
	protos := make([]*pb.ChatMessage, len(entities))
	for i := range entities {
		protos[i] = MessageToProto(&entities[i])
	}
	return protos
}
