// internal/visitor/grpc/handler.go
package grpc

import (
	"context"
	"errors"
	"dan-ai/internal/visitor/mapper"
	"dan-ai/internal/visitor/service"
	pb "dan-ai/proto/visitor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Handler struct {
	pb.UnimplementedVisitorServiceServer
	svc service.Service
}

// NewHandler creates a new gRPC handler for the Visitor service.
func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterVisitor(ctx context.Context, req *pb.RegisterVisitorRequest) (*pb.RegisterVisitorResponse, error) {
	v, err := h.svc.Register(ctx, req.GetVisitorId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register visitor: %v", err)
	}

	return &pb.RegisterVisitorResponse{
		Visitor: mapper.ToProto(v),
	}, nil
}

func (h *Handler) GetVisitor(ctx context.Context, req *pb.GetVisitorRequest) (*pb.GetVisitorResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	v, err := h.svc.Get(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "visitor not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get visitor: %v", err)
	}

	return &pb.GetVisitorResponse{
		Visitor: mapper.ToProto(v),
	}, nil
}
