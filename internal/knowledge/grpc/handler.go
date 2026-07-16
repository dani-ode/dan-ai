package grpc

import (
	"context"

	"dan-ai/internal/knowledge/mapper"
	"dan-ai/internal/knowledge/service"
	pb "dan-ai/proto/knowledge"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KnowledgeHandler struct {
	pb.UnimplementedKnowledgeServiceServer
	svc service.Service
}

func NewKnowledgeHandler(svc service.Service) *KnowledgeHandler {
	return &KnowledgeHandler{svc: svc}
}

func (h *KnowledgeHandler) GetDocument(ctx context.Context, req *pb.GetDocumentRequest) (*pb.KnowledgeDocument, error) {
	doc, err := h.svc.GetDocument(ctx, req.Id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Error(codes.NotFound, "document not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get document: %v", err)
	}

	return mapper.DocumentEntityToProto(doc), nil
}

func (h *KnowledgeHandler) ListDocuments(ctx context.Context, req *pb.ListDocumentsRequest) (*pb.ListDocumentsResponse, error) {
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	docs, total, err := h.svc.ListDocuments(ctx, page, pageSize, req.SourceType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list documents: %v", err)
	}

	var pbDocs []*pb.KnowledgeDocument
	for _, doc := range docs {
		// need local variable to take pointer
		d := doc
		pbDocs = append(pbDocs, mapper.DocumentEntityToProto(&d))
	}

	return &pb.ListDocumentsResponse{
		Documents:  pbDocs,
		TotalCount: int32(total),
	}, nil
}

func (h *KnowledgeHandler) ListChunks(ctx context.Context, req *pb.ListChunksRequest) (*pb.ListChunksResponse, error) {
	chunks, err := h.svc.ListChunks(ctx, req.DocumentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list chunks: %v", err)
	}

	var pbChunks []*pb.KnowledgeChunk
	for _, chunk := range chunks {
		c := chunk
		pbChunks = append(pbChunks, mapper.ChunkEntityToProto(&c))
	}

	return &pb.ListChunksResponse{
		Chunks: pbChunks,
	}, nil
}
