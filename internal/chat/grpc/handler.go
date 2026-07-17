// internal/chat/grpc/handler.go
package grpc

import (
	"context"
	"errors"
	"dan-ai/internal/chat/mapper"
	"dan-ai/internal/chat/service"
	visitorsvc "dan-ai/internal/visitor/service"
	pb "dan-ai/proto/chat"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Handler struct {
	pb.UnimplementedChatServiceServer
	svc        service.Service
	visitorSvc visitorsvc.Service
}

// NewHandler creates a new gRPC handler for the Chat service.
func NewHandler(svc service.Service, visitorSvc visitorsvc.Service) *Handler {
	return &Handler{
		svc:        svc,
		visitorSvc: visitorSvc,
	}
}

// --- Session handlers ---

func (h *Handler) CreateSession(ctx context.Context, req *pb.CreateSessionRequest) (*pb.CreateSessionResponse, error) {
	if req.GetVisitorId() == "" {
		return nil, status.Error(codes.InvalidArgument, "visitor_id is required")
	}

	session, err := h.svc.CreateSession(ctx, req.GetVisitorId(), req.GetPromptId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	return &pb.CreateSessionResponse{
		Session: mapper.SessionToProto(session),
	}, nil
}

func (h *Handler) ListSessions(ctx context.Context, req *pb.ListSessionsRequest) (*pb.ListSessionsResponse, error) {
	if req.GetVisitorId() == "" {
		return nil, status.Error(codes.InvalidArgument, "visitor_id is required")
	}

	sessions, err := h.svc.ListSessions(ctx, req.GetVisitorId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list sessions: %v", err)
	}

	return &pb.ListSessionsResponse{
		Sessions: mapper.SessionsToProto(sessions),
	}, nil
}

func (h *Handler) RenameSession(ctx context.Context, req *pb.RenameSessionRequest) (*pb.RenameSessionResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	session, err := h.svc.RenameSession(ctx, req.GetId(), req.GetTitle())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "session not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to rename session: %v", err)
	}

	return &pb.RenameSessionResponse{
		Session: mapper.SessionToProto(session),
	}, nil
}

func (h *Handler) DeleteSession(ctx context.Context, req *pb.DeleteSessionRequest) (*pb.DeleteSessionResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := h.svc.DeleteSession(ctx, req.GetId()); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "session not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete session: %v", err)
	}

	return &pb.DeleteSessionResponse{
		Success: true,
	}, nil
}

// --- Message handlers ---

func (h *Handler) CreateMessage(ctx context.Context, req *pb.CreateMessageRequest) (*pb.CreateMessageResponse, error) {
	if req.GetSessionId() == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	if req.GetRole() == "" {
		return nil, status.Error(codes.InvalidArgument, "role is required")
	}
	if req.GetContent() == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}

	msg, err := h.svc.CreateMessage(ctx, req.GetSessionId(), req.GetRole(), req.GetContent())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create message: %v", err)
	}

	return &pb.CreateMessageResponse{
		Message: mapper.MessageToProto(msg),
	}, nil
}

func (h *Handler) ListMessages(ctx context.Context, req *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
	if req.GetSessionId() == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	messages, err := h.svc.ListMessages(ctx, req.GetSessionId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list messages: %v", err)
	}

	return &pb.ListMessagesResponse{
		Messages: mapper.MessagesToProto(messages),
	}, nil
}

func (h *Handler) DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := h.svc.DeleteMessage(ctx, req.GetId()); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "message not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete message: %v", err)
	}

	return &pb.DeleteMessageResponse{
		Success: true,
	}, nil
}

// --- Unified Chat handlers ---

func (h *Handler) SendChatMessage(ctx context.Context, req *pb.SendChatMessageRequest) (*pb.SendChatMessageResponse, error) {
	assistantMsg, err := h.svc.SendChatMessage(ctx, req.GetSessionId(), req.GetVisitorId(), req.GetPromptId(), req.GetContent())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send chat message: %v", err)
	}

	// Fetch recent messages for the response (the last 3 messages)
	messages, err := h.svc.ListMessages(ctx, assistantMsg.SessionID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.Internal, "failed to list messages: %v", err)
	}

	// Slicing the last 3 messages (e.g. recent dialogue context)
	startIndex := 0
	if len(messages) > 3 {
		startIndex = len(messages) - 3
	}
	recentMessages := messages[startIndex:]

	var recentContent []*pb.RecentMessage
	for _, m := range recentMessages {
		recentContent = append(recentContent, &pb.RecentMessage{
			Role:      m.Role,
			Content:   m.Content,
			CreatedAt: timestamppb.New(m.CreatedAt),
		})
	}

	session, err := h.svc.GetSession(ctx, assistantMsg.SessionID)
	visitorID := req.GetVisitorId()
	if err == nil {
		visitorID = session.VisitorID
	}

	return &pb.SendChatMessageResponse{
		SessionId:      assistantMsg.SessionID,
		VisitorId:      visitorID,
		Content:        assistantMsg.Content,
		RecentMessages: recentContent,
	}, nil
}
