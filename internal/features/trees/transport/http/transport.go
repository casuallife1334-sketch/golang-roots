package http

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/server"
	"net/http"
)

type TreesService interface {
	CreateTree(context.Context, string, domain.CreateTreeInput) (domain.Tree, error)
	GetTrees(context.Context, string) ([]domain.Tree, error)
	GetTree(context.Context, string, string) (domain.Tree, error)
	PatchTree(context.Context, string, string, domain.PatchTreeInput) (domain.Tree, error)
	DeleteTree(context.Context, string, string) error
}

type TreesHTTPHandler struct {
	treesService TreesService
}

func NewTreesHTTPHandler(treesService TreesService) *TreesHTTPHandler {
	return &TreesHTTPHandler{treesService: treesService}
}

func (h *TreesHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{Method: http.MethodPost, Path: "/trees", Handler: h.CreateTree},
		{Method: http.MethodGet, Path: "/trees", Handler: h.GetTrees},
		{Method: http.MethodGet, Path: "/trees/{tree_id}", Handler: h.GetTree},
		{Method: http.MethodPatch, Path: "/trees/{tree_id}", Handler: h.PatchTree},
		{Method: http.MethodDelete, Path: "/trees/{tree_id}", Handler: h.DeleteTree},
	}
}
