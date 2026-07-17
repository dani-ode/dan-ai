// internal/chat/grpc/handler_test.go
package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"dan-ai/internal/chat/entity"
	visitorentity "dan-ai/internal/visitor/entity"
	visitorsvc "dan-ai/internal/visitor/service"
	pb "dan-ai/proto/chat"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// Manual mocks
// ---------------------------------------------------------------------------

// mockChatService implements service.Service.
type mockChatService struct {
	getSessionFn      func(ctx context.Context, id string) (*entity.ChatSession, error)
	createSessionFn   func(ctx context.Context, visitorID, promptID string) (*entity.ChatSession, error)
	listMessagesFn    func(ctx context.Context, sessionID string) ([]entity.ChatMessage, error)
	createMessageFn   func(ctx context.Context, sessionID, role, content string) (*entity.ChatMessage, error)
	sendChatMessageFn func(ctx context.Context, sessionID, visitorID, promptID, content string) (*entity.ChatMessage, error)
}

func (m *mockChatService) GetSession(ctx context.Context, id string) (*entity.ChatSession, error) {
	if m.getSessionFn != nil {
		return m.getSessionFn(ctx, id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockChatService) CreateSession(ctx context.Context, visitorID, promptID string) (*entity.ChatSession, error) {
	if m.createSessionFn != nil {
		return m.createSessionFn(ctx, visitorID, promptID)
	}
	return &entity.ChatSession{ID: "new-session-id", VisitorID: visitorID, PromptID: promptID}, nil
}

func (m *mockChatService) ListMessages(ctx context.Context, sessionID string) ([]entity.ChatMessage, error) {
	if m.listMessagesFn != nil {
		return m.listMessagesFn(ctx, sessionID)
	}
	return nil, nil
}

func (m *mockChatService) CreateMessage(ctx context.Context, sessionID, role, content string) (*entity.ChatMessage, error) {
	if m.createMessageFn != nil {
		return m.createMessageFn(ctx, sessionID, role, content)
	}
	return &entity.ChatMessage{ID: "msg-id", SessionID: sessionID, Role: role, Content: content}, nil
}

func (m *mockChatService) ListSessions(ctx context.Context, visitorID string) ([]entity.ChatSession, error) {
	return nil, nil
}
func (m *mockChatService) RenameSession(ctx context.Context, id, title string) (*entity.ChatSession, error) {
	return nil, nil
}
func (m *mockChatService) DeleteSession(ctx context.Context, id string) error { return nil }
func (m *mockChatService) DeleteMessage(ctx context.Context, id string) error { return nil }
func (m *mockChatService) SendChatMessage(ctx context.Context, sessionID, visitorID, promptID, content string) (*entity.ChatMessage, error) {
	if m.sendChatMessageFn != nil {
		return m.sendChatMessageFn(ctx, sessionID, visitorID, promptID, content)
	}
	return &entity.ChatMessage{ID: "msg-id", SessionID: sessionID, Role: "assistant", Content: "reply-content"}, nil
}

// mockVisitorService implements visitorsvc.Service.
type mockVisitorService struct {
	registerFn func(ctx context.Context, visitorID string) (*visitorentity.Visitor, error)
}

func (m *mockVisitorService) Register(ctx context.Context, visitorID string) (*visitorentity.Visitor, error) {
	if m.registerFn != nil {
		return m.registerFn(ctx, visitorID)
	}
	return &visitorentity.Visitor{ID: "visitor-id"}, nil
}

func (m *mockVisitorService) Get(ctx context.Context, id string) (*visitorentity.Visitor, error) {
	return nil, nil
}

// compile-time interface checks
var _ visitorsvc.Service = (*mockVisitorService)(nil)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newTestHandler(svc *mockChatService, vsvc visitorsvc.Service) *Handler {
	return NewHandler(svc, vsvc)
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// TestSendChatMessage_ExistingSession verifies that SendChatMessage is successfully forwarded to the service.
func TestSendChatMessage_ExistingSession(t *testing.T) {
	ctx := context.Background()

	svc := &mockChatService{
		sendChatMessageFn: func(ctx context.Context, sessionID, visitorID, promptID, content string) (*entity.ChatMessage, error) {
			assert.Equal(t, "session-123", sessionID)
			assert.Equal(t, "new message", content)
			return &entity.ChatMessage{
				ID:        "msg-3",
				SessionID: "session-123",
				Role:      "assistant",
				Content:   "reply from assistant",
			}, nil
		},
		listMessagesFn: func(_ context.Context, sessionID string) ([]entity.ChatMessage, error) {
			assert.Equal(t, "session-123", sessionID)
			return []entity.ChatMessage{
				{ID: "m1", Role: "user", Content: "hello", CreatedAt: time.Now()},
				{ID: "m2", Role: "assistant", Content: "hi", CreatedAt: time.Now()},
				{ID: "msg-3", Role: "assistant", Content: "reply from assistant", CreatedAt: time.Now()},
			}, nil
		},
		getSessionFn: func(_ context.Context, id string) (*entity.ChatSession, error) {
			return &entity.ChatSession{ID: id, VisitorID: "visitor-456"}, nil
		},
	}

	h := newTestHandler(svc, &mockVisitorService{})
	resp, err := h.SendChatMessage(ctx, &pb.SendChatMessageRequest{
		SessionId: "session-123",
		Content:   "new message",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "session-123", resp.GetSessionId())
	assert.Equal(t, "visitor-456", resp.GetVisitorId())
	assert.Equal(t, "reply from assistant", resp.GetContent())
	assert.Len(t, resp.GetRecentMessages(), 3)
}

// TestSendChatMessage_ExistingSession_InheritsVisitorAndPrompt verifies that
// if visitor_id and prompt_id are empty in the request, they are filled from
// the existing session.
func TestSendChatMessage_ExistingSession_InheritsVisitorAndPrompt(t *testing.T) {
	ctx := context.Background()

	svc := &mockChatService{
		sendChatMessageFn: func(ctx context.Context, sessionID, visitorID, promptID, content string) (*entity.ChatMessage, error) {
			return &entity.ChatMessage{
				ID:        "msg-3",
				SessionID: "sess-1",
				Role:      "assistant",
				Content:   "reply from assistant",
			}, nil
		},
		getSessionFn: func(_ context.Context, _ string) (*entity.ChatSession, error) {
			return &entity.ChatSession{
				ID:        "sess-1",
				VisitorID: "vis-from-session",
				PromptID:  "prompt-from-session",
			}, nil
		},
	}

	h := newTestHandler(svc, &mockVisitorService{})
	resp, err := h.SendChatMessage(ctx, &pb.SendChatMessageRequest{
		SessionId: "sess-1",
		Content:   "hi",
	})

	require.NoError(t, err)
	assert.Equal(t, "vis-from-session", resp.GetVisitorId())
}

// TestSendChatMessage_NewSession_NoSessionID verifies that when no session_id
// is provided, the handler delegates new session creation to service.
func TestSendChatMessage_NewSession_NoSessionID(t *testing.T) {
	ctx := context.Background()

	svc := &mockChatService{
		sendChatMessageFn: func(ctx context.Context, sessionID, visitorID, promptID, content string) (*entity.ChatMessage, error) {
			assert.Equal(t, "", sessionID)
			assert.Equal(t, "new-visitor-id", visitorID)
			return &entity.ChatMessage{
				ID:        "msg-1",
				SessionID: "new-sess-id",
				Role:      "assistant",
				Content:   "reply from assistant",
			}, nil
		},
		getSessionFn: func(_ context.Context, id string) (*entity.ChatSession, error) {
			return &entity.ChatSession{ID: id, VisitorID: "new-visitor-id"}, nil
		},
	}

	h := newTestHandler(svc, &mockVisitorService{})
	resp, err := h.SendChatMessage(ctx, &pb.SendChatMessageRequest{
		VisitorId: "new-visitor-id",
		Content:   "hello world",
	})

	require.NoError(t, err)
	assert.Equal(t, "new-sess-id", resp.GetSessionId())
	assert.Equal(t, "new-visitor-id", resp.GetVisitorId())
	assert.Equal(t, "reply from assistant", resp.GetContent())
}

// TestSendChatMessage_SessionNotFound_CreatesNew verifies that when session_id
// is provided but GetSession returns ErrRecordNotFound, a new session is
// created.
func TestSendChatMessage_SessionNotFound_CreatesNew(t *testing.T) {
	ctx := context.Background()

	svc := &mockChatService{
		sendChatMessageFn: func(ctx context.Context, sessionID, visitorID, promptID, content string) (*entity.ChatMessage, error) {
			return &entity.ChatMessage{ID: "msg-1", SessionID: "created-sess"}, nil
		},
		getSessionFn: func(_ context.Context, _ string) (*entity.ChatSession, error) {
			return &entity.ChatSession{ID: "created-sess", VisitorID: "registered-visitor"}, nil
		},
	}

	h := newTestHandler(svc, &mockVisitorService{})
	resp, err := h.SendChatMessage(ctx, &pb.SendChatMessageRequest{
		SessionId: "ghost-session",
		Content:   "anyone there?",
	})

	require.NoError(t, err)
	assert.Equal(t, "created-sess", resp.GetSessionId())
}

// TestSendChatMessage_RecentMessages_SlicesLastThree verifies that only the
// last 3 messages (out of more than 3 prior messages) are returned.
func TestSendChatMessage_RecentMessages_SlicesLastThree(t *testing.T) {
	ctx := context.Background()

	msgs := []entity.ChatMessage{
		{ID: "m1", Role: "user", Content: "msg1", CreatedAt: time.Now()},
		{ID: "m2", Role: "assistant", Content: "msg2", CreatedAt: time.Now()},
		{ID: "m3", Role: "user", Content: "msg3", CreatedAt: time.Now()},
		{ID: "m4", Role: "assistant", Content: "msg4", CreatedAt: time.Now()},
		{ID: "m5", Role: "user", Content: "msg5", CreatedAt: time.Now()},
	}

	svc := &mockChatService{
		sendChatMessageFn: func(ctx context.Context, sessionID, visitorID, promptID, content string) (*entity.ChatMessage, error) {
			return &entity.ChatMessage{ID: "latest", SessionID: "sess"}, nil
		},
		getSessionFn: func(_ context.Context, _ string) (*entity.ChatSession, error) {
			return &entity.ChatSession{ID: "sess", VisitorID: "vis"}, nil
		},
		listMessagesFn: func(_ context.Context, _ string) ([]entity.ChatMessage, error) {
			return msgs, nil
		},
	}

	h := newTestHandler(svc, &mockVisitorService{})
	resp, err := h.SendChatMessage(ctx, &pb.SendChatMessageRequest{
		SessionId: "sess",
		Content:   "latest message",
	})

	require.NoError(t, err)
	assert.Len(t, resp.GetRecentMessages(), 3)
	assert.Equal(t, "msg3", resp.GetRecentMessages()[0].GetContent())
	assert.Equal(t, "msg4", resp.GetRecentMessages()[1].GetContent())
	assert.Equal(t, "msg5", resp.GetRecentMessages()[2].GetContent())
}

// TestSendChatMessage_GetSession_InternalError verifies that an unexpected
// error from SendChatMessage propagates as an Internal gRPC status.
func TestSendChatMessage_GetSession_InternalError(t *testing.T) {
	ctx := context.Background()

	svc := &mockChatService{
		sendChatMessageFn: func(ctx context.Context, sessionID, visitorID, promptID, content string) (*entity.ChatMessage, error) {
			return nil, errors.New("db timeout")
		},
	}

	h := newTestHandler(svc, &mockVisitorService{})
	_, err := h.SendChatMessage(ctx, &pb.SendChatMessageRequest{
		SessionId: "sess-bad",
		Content:   "hello",
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "failed to send chat message")
}

// TestSendChatMessage_ListMessages_InternalError verifies that a non-NotFound
// error from ListMessages returns an Internal gRPC error.
func TestSendChatMessage_ListMessages_InternalError(t *testing.T) {
	ctx := context.Background()

	svc := &mockChatService{
		sendChatMessageFn: func(ctx context.Context, sessionID, visitorID, promptID, content string) (*entity.ChatMessage, error) {
			return &entity.ChatMessage{ID: "latest", SessionID: "sess"}, nil
		},
		listMessagesFn: func(_ context.Context, _ string) ([]entity.ChatMessage, error) {
			return nil, errors.New("query failed")
		},
	}

	h := newTestHandler(svc, &mockVisitorService{})
	_, err := h.SendChatMessage(ctx, &pb.SendChatMessageRequest{
		SessionId: "sess",
		Content:   "hi",
	})

	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "failed to list messages")
}

