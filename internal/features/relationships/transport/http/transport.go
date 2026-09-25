package http

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/server"
	"net/http"
)

type RelationshipsService interface {
	CreateRelationship(context.Context, string, string, domain.CreateRelationshipInput) (domain.Relationship, error)
	GetRelationship(context.Context, string, string, string) (domain.Relationship, error)
	GetRelationships(context.Context, string, string, string) ([]domain.Relationship, error)
	PatchRelationship(context.Context, string, string, string, domain.PatchRelationshipInput) (domain.Relationship, error)
	DeleteRelationship(context.Context, string, string, string) error
}

type RelationshipsHTTPHandler struct {
	relationshipsService RelationshipsService
}

func NewRelationshipsHTTPHandlers(relationshipsService RelationshipsService) *RelationshipsHTTPHandler {
	return &RelationshipsHTTPHandler{relationshipsService: relationshipsService}
}

func (h *RelationshipsHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{Method: http.MethodPost, Path: "/trees/{tree_id}/relationships", Handler: h.CreateRelationship},
		{Method: http.MethodGet, Path: "/trees/{tree_id}/relationships", Handler: h.GetRelationships},
		{Method: http.MethodGet, Path: "/trees/{tree_id}/relationships/{id}", Handler: h.GetRelationship},
		{Method: http.MethodPatch, Path: "/trees/{tree_id}/relationships/{id}", Handler: h.PatchRelationship},
		{Method: http.MethodDelete, Path: "/trees/{tree_id}/relationships/{id}", Handler: h.DeleteRelationship},
	}
}
