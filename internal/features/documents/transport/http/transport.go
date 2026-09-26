package http

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/server"
	"genealogy-tree/internal/features/documents/service"
	"io"
	"net/http"
)

type DocumentsService interface {
	CreateDocument(context.Context, string, string, service.UploadInput) (domain.Document, error)
	GetDocuments(context.Context, string, string, domain.DocumentOwner) ([]domain.Document, error)
	OpenDocument(context.Context, string, string, string) (domain.Document, io.ReadCloser, error)
	DeleteDocument(context.Context, string, string, string) error
}

type DocumentsHTTPHandler struct {
	documentsService DocumentsService
}

func NewDocumentsHTTPHandler(documentsService DocumentsService) *DocumentsHTTPHandler {
	return &DocumentsHTTPHandler{documentsService: documentsService}
}

func (h *DocumentsHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{Method: http.MethodPost, Path: "/trees/{tree_id}/documents", Handler: h.CreateDocument},
		{Method: http.MethodGet, Path: "/trees/{tree_id}/documents", Handler: h.GetDocuments},
		{Method: http.MethodGet, Path: "/trees/{tree_id}/documents/{id}", Handler: h.GetDocument},
		{Method: http.MethodDelete, Path: "/trees/{tree_id}/documents/{id}", Handler: h.DeleteDocument},
	}
}
