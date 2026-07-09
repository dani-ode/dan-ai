// internal/aimodel/grpc/handler.go
package grpc

import (
	"context"
	"errors"
	"portfolio-ai/internal/aimodel/mapper"
	"portfolio-ai/internal/aimodel/service"
	pb "portfolio-ai/proto/aimodel"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Handler struct {
	pb.UnimplementedAIModelServiceServer
	svc service.Service
}

// NewHandler creates a new gRPC handler for the AIModel service.
func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateAIModel(ctx context.Context, req *pb.CreateAIModelRequest) (*pb.CreateAIModelResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.GetProvider() == "" {
		return nil, status.Error(codes.InvalidArgument, "provider is required")
	}

	m := mapper.ToEntityFromCreate(req)
	if err := h.svc.Create(ctx, m); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create ai model: %v", err)
	}

	return &pb.CreateAIModelResponse{
		AiModel: mapper.ToProto(m),
	}, nil
}

func (h *Handler) GetAIModel(ctx context.Context, req *pb.GetAIModelRequest) (*pb.GetAIModelResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	m, err := h.svc.Get(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "ai model not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get ai model: %v", err)
	}

	return &pb.GetAIModelResponse{
		AiModel: mapper.ToProto(m),
	}, nil
}

func (h *Handler) ListAIModels(ctx context.Context, req *pb.ListAIModelsRequest) (*pb.ListAIModelsResponse, error) {
	models, err := h.svc.List(ctx, req.GetEnabledOnly())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list ai models: %v", err)
	}

	return &pb.ListAIModelsResponse{
		AiModels: mapper.ToProtoList(models),
	}, nil
}

func (h *Handler) UpdateAIModel(ctx context.Context, req *pb.UpdateAIModelRequest) (*pb.UpdateAIModelResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.GetProvider() == "" {
		return nil, status.Error(codes.InvalidArgument, "provider is required")
	}

	m := mapper.ToEntityFromUpdate(req)
	if err := h.svc.Update(ctx, m); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "ai model not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update ai model: %v", err)
	}

	return &pb.UpdateAIModelResponse{
		AiModel: mapper.ToProto(m),
	}, nil
}

func (h *Handler) DeleteAIModel(ctx context.Context, req *pb.DeleteAIModelRequest) (*pb.DeleteAIModelResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := h.svc.Delete(ctx, req.GetId()); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "ai model not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete ai model: %v", err)
	}

	return &pb.DeleteAIModelResponse{
		Success: true,
	}, nil
}
