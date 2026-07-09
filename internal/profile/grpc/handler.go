// internal/profile/grpc/handler.go
package grpc

import (
	"context"
	"errors"
	"portfolio-ai/internal/profile/mapper"
	"portfolio-ai/internal/profile/service"
	pb "portfolio-ai/proto/profile"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Handler struct {
	pb.UnimplementedProfileServiceServer
	svc service.Service
}

// NewHandler creates a new gRPC handler for the Profile service.
func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	p, err := h.svc.Get(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "profile not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get profile: %v", err)
	}

	return &pb.GetProfileResponse{
		Profile: mapper.ToProto(p),
	}, nil
}

func (h *Handler) CreateProfile(ctx context.Context, req *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	if req.GetFullName() == "" {
		return nil, status.Error(codes.InvalidArgument, "full_name is required")
	}

	p := mapper.ToEntityFromCreate(req)
	if err := h.svc.Create(ctx, p); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create profile: %v", err)
	}

	return &pb.CreateProfileResponse{
		Profile: mapper.ToProto(p),
	}, nil
}

func (h *Handler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.GetFullName() == "" {
		return nil, status.Error(codes.InvalidArgument, "full_name is required")
	}

	p := mapper.ToEntityFromUpdate(req)
	if err := h.svc.Update(ctx, p); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "profile not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update profile: %v", err)
	}

	return &pb.UpdateProfileResponse{
		Profile: mapper.ToProto(p),
	}, nil
}

func (h *Handler) DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest) (*pb.DeleteProfileResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := h.svc.Delete(ctx, req.GetId()); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "profile not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete profile: %v", err)
	}

	return &pb.DeleteProfileResponse{
		Success: true,
	}, nil
}
