package http

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/server"
	"net/http"
)

type RelationshipsService interface {
	CreateRelationship(context.Context, domain.CreateRelationshipInput) (domain.Relationship, error)
	GetRelationship(context.Context, string) (domain.Relationship, error)
	DeleteRelationship(context.Context, string) error
}

type RelationshipsHTTPHandler struct {
	relationshipsService RelationshipsService
}

func NewRelationshipsHTTPHandlers(relationshipsService RelationshipsService) *RelationshipsHTTPHandler {
	return &RelationshipsHTTPHandler{relationshipsService: relationshipsService}
}

func (h *RelationshipsHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{Method: http.MethodPost, Path: "/relationships", Handler: h.CreateRelationship},
		{Method: http.MethodGet, Path: "/relationships/{id}", Handler: h.GetRelationship},
		{Method: http.MethodDelete, Path: "/relationships/{id}", Handler: h.DeleteRelationship},
	}
}
