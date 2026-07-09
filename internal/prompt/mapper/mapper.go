// internal/prompt/mapper/mapper.go
package mapper

import (
	"portfolio-ai/internal/prompt/entity"
	pb "portfolio-ai/proto/prompt"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProto maps a domain Prompt entity to a protobuf Prompt message.
func ToProto(e *entity.Prompt) *pb.Prompt {
	if e == nil {
		return nil
	}
	return &pb.Prompt{
		Id:           e.ID,
		Name:         e.Name,
		SystemPrompt: e.SystemPrompt,
		Description:  e.Description,
		ModelId:      e.ModelID,
		Active:       e.Active,
		Version:      e.Version,
		CreatedAt:    timestamppb.New(e.CreatedAt),
	}
}

// ToProtoList maps a slice of Prompt entities to protobuf Prompt messages.
func ToProtoList(entities []entity.Prompt) []*pb.Prompt {
	protos := make([]*pb.Prompt, len(entities))
	for i := range entities {
		protos[i] = ToProto(&entities[i])
	}
	return protos
}

// ToEntityFromCreate maps a CreatePromptRequest to a domain Prompt entity.
func ToEntityFromCreate(r *pb.CreatePromptRequest) *entity.Prompt {
	if r == nil {
		return nil
	}
	return &entity.Prompt{
		Name:         r.Name,
		SystemPrompt: r.SystemPrompt,
		Description:  r.Description,
		ModelID:      r.ModelId,
		Active:       r.Active,
	}
}

// ToEntityFromUpdate maps an UpdatePromptRequest to a domain Prompt entity.
func ToEntityFromUpdate(r *pb.UpdatePromptRequest) *entity.Prompt {
	if r == nil {
		return nil
	}
	return &entity.Prompt{
		ID:           r.Id,
		Name:         r.Name,
		SystemPrompt: r.SystemPrompt,
		Description:  r.Description,
		ModelID:      r.ModelId,
		Active:       r.Active,
	}
}
