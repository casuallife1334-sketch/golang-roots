package http

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/server"
	"net/http"
)

type RelationshipsService interface {
	GetRelationship(context.Context, string, string, string) (domain.Relationship, error)
	GetRelationships(context.Context, string, string, string) ([]domain.Relationship, error)
	PatchRelationship(context.Context, string, string, string, domain.PatchRelationshipInput) (domain.Relationship, error)
	DeleteRelationship(context.Context, string, string, string) error
}

type RelationshipCreator interface {
	CreateRelationship(context.Context, string, string, domain.CreateRelationshipCommand) (domain.Relationship, error)
}

type RelationshipsHTTPHandler struct {
	relationshipsService RelationshipsService
	relationshipCreator  RelationshipCreator
}

func NewRelationshipsHTTPHandlers(relationshipsService RelationshipsService, relationshipCreator RelationshipCreator) *RelationshipsHTTPHandler {
	return &RelationshipsHTTPHandler{
		relationshipsService: relationshipsService,
		relationshipCreator:  relationshipCreator,
	}
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
