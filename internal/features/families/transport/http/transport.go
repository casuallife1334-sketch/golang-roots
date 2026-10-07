package http

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/server"
	"net/http"
)

type FamiliesService interface {
	CreateFamily(context.Context, string, string, domain.CreateFamilyInput) (domain.Family, error)
	GetFamilies(context.Context, string, string) ([]domain.Family, error)
	GetFamily(context.Context, string, string, string) (domain.Family, error)
	PatchFamily(context.Context, string, string, string, domain.PatchFamilyInput) (domain.Family, error)
	DeleteFamily(context.Context, string, string, string) error
	AddMember(context.Context, string, string, string, domain.AddFamilyMemberInput) error
	RemoveMember(context.Context, string, string, string, string) error
	AttachRelationship(context.Context, string, string, string, string) error
	DetachRelationship(context.Context, string, string, string, string) error
}

type FamiliesHTTPHandler struct {
	familiesService FamiliesService
}

func NewFamiliesHTTPHandler(familiesService FamiliesService) *FamiliesHTTPHandler {
	return &FamiliesHTTPHandler{familiesService: familiesService}
}

func (h *FamiliesHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{Method: http.MethodPost, Path: "/trees/{tree_id}/families", Handler: h.CreateFamily},
		{Method: http.MethodGet, Path: "/trees/{tree_id}/families", Handler: h.GetFamilies},
		{Method: http.MethodGet, Path: "/trees/{tree_id}/families/{family_id}", Handler: h.GetFamily},
		{Method: http.MethodPatch, Path: "/trees/{tree_id}/families/{family_id}", Handler: h.PatchFamily},
		{Method: http.MethodDelete, Path: "/trees/{tree_id}/families/{family_id}", Handler: h.DeleteFamily},
		{Method: http.MethodPost, Path: "/trees/{tree_id}/families/{family_id}/members", Handler: h.AddMember},
		{Method: http.MethodDelete, Path: "/trees/{tree_id}/families/{family_id}/members/{person_id}", Handler: h.RemoveMember},
		{Method: http.MethodPost, Path: "/trees/{tree_id}/families/{family_id}/relationships", Handler: h.AttachRelationship},
		{Method: http.MethodDelete, Path: "/trees/{tree_id}/families/{family_id}/relationships/{relationship_id}", Handler: h.DetachRelationship},
	}
}
