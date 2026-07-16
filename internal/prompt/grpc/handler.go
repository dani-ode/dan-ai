// internal/prompt/grpc/handler.go
package grpc

import (
	"context"
	"errors"
	"dan-ai/internal/prompt/mapper"
	"dan-ai/internal/prompt/service"
	pb "dan-ai/proto/prompt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Handler struct {
	pb.UnimplementedPromptServiceServer
	svc service.Service
}

// NewHandler creates a new gRPC handler for the Prompt service.
func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreatePrompt(ctx context.Context, req *pb.CreatePromptRequest) (*pb.CreatePromptResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.GetSystemPrompt() == "" {
		return nil, status.Error(codes.InvalidArgument, "system_prompt is required")
	}

	p := mapper.ToEntityFromCreate(req)
	if err := h.svc.Create(ctx, p); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create prompt: %v", err)
	}

	return &pb.CreatePromptResponse{
		Prompt: mapper.ToProto(p),
	}, nil
}

func (h *Handler) GetPrompt(ctx context.Context, req *pb.GetPromptRequest) (*pb.GetPromptResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	p, err := h.svc.Get(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "prompt not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get prompt: %v", err)
	}

	return &pb.GetPromptResponse{
		Prompt: mapper.ToProto(p),
	}, nil
}

func (h *Handler) ListPrompts(ctx context.Context, req *pb.ListPromptsRequest) (*pb.ListPromptsResponse, error) {
	prompts, err := h.svc.List(ctx, req.GetActiveOnly())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list prompts: %v", err)
	}

	return &pb.ListPromptsResponse{
		Prompts: mapper.ToProtoList(prompts),
	}, nil
}

func (h *Handler) UpdatePrompt(ctx context.Context, req *pb.UpdatePromptRequest) (*pb.UpdatePromptResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.GetSystemPrompt() == "" {
		return nil, status.Error(codes.InvalidArgument, "system_prompt is required")
	}

	p := mapper.ToEntityFromUpdate(req)
	if err := h.svc.Update(ctx, p); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "prompt not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update prompt: %v", err)
	}

	return &pb.UpdatePromptResponse{
		Prompt: mapper.ToProto(p),
	}, nil
}

func (h *Handler) DeletePrompt(ctx context.Context, req *pb.DeletePromptRequest) (*pb.DeletePromptResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := h.svc.Delete(ctx, req.GetId()); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "prompt not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete prompt: %v", err)
	}

	return &pb.DeletePromptResponse{
		Success: true,
	}, nil
}
