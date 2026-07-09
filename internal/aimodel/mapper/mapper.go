// internal/aimodel/mapper/mapper.go
package mapper

import (
	"portfolio-ai/internal/aimodel/entity"
	pb "portfolio-ai/proto/aimodel"
)

// ToProto maps a domain AIModel entity to a protobuf AIModel message.
func ToProto(e *entity.AIModel) *pb.AIModel {
	if e == nil {
		return nil
	}
	return &pb.AIModel{
		Id:             e.ID,
		Name:           e.Name,
		Provider:       e.Provider,
		Temperature:    e.Temperature,
		MaxTokens:      e.MaxTokens,
		ContextWindow:  e.ContextWindow,
		SupportsTools:  e.SupportsTools,
		SupportsStream: e.SupportsStream,
		Enabled:        e.Enabled,
	}
}

// ToProtoList maps a slice of AIModel entities to protobuf messages.
func ToProtoList(entities []entity.AIModel) []*pb.AIModel {
	protos := make([]*pb.AIModel, len(entities))
	for i := range entities {
		protos[i] = ToProto(&entities[i])
	}
	return protos
}

// ToEntityFromCreate maps a CreateAIModelRequest to a domain AIModel entity.
func ToEntityFromCreate(r *pb.CreateAIModelRequest) *entity.AIModel {
	if r == nil {
		return nil
	}
	return &entity.AIModel{
		Name:           r.Name,
		Provider:       r.Provider,
		Temperature:    r.Temperature,
		MaxTokens:      r.MaxTokens,
		ContextWindow:  r.ContextWindow,
		SupportsTools:  r.SupportsTools,
		SupportsStream: r.SupportsStream,
		Enabled:        r.Enabled,
	}
}

// ToEntityFromUpdate maps an UpdateAIModelRequest to a domain AIModel entity.
func ToEntityFromUpdate(r *pb.UpdateAIModelRequest) *entity.AIModel {
	if r == nil {
		return nil
	}
	return &entity.AIModel{
		ID:             r.Id,
		Name:           r.Name,
		Provider:       r.Provider,
		Temperature:    r.Temperature,
		MaxTokens:      r.MaxTokens,
		ContextWindow:  r.ContextWindow,
		SupportsTools:  r.SupportsTools,
		SupportsStream: r.SupportsStream,
		Enabled:        r.Enabled,
	}
}
