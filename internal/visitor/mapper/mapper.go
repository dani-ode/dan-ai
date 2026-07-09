// internal/visitor/mapper/mapper.go
package mapper

import (
	"portfolio-ai/internal/visitor/entity"
	pb "portfolio-ai/proto/visitor"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProto maps a domain Visitor entity to a protobuf Visitor message.
func ToProto(e *entity.Visitor) *pb.Visitor {
	if e == nil {
		return nil
	}
	return &pb.Visitor{
		Id:            e.ID,
		FirstSeenAt:   timestamppb.New(e.FirstSeenAt),
		LastSeenAt:    timestamppb.New(e.LastSeenAt),
		TotalMessages: e.TotalMessages,
		CreatedAt:     timestamppb.New(e.CreatedAt),
		UpdatedAt:     timestamppb.New(e.UpdatedAt),
	}
}
