package http

import (
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// PatchTree godoc
// @Summary Изменение дерева
// @Description Изменение названия дерева владельцем
// @Tags trees
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tree_id path string true "ULID дерева"
// @Param request body PatchTreeRequest true "PatchTree тело запроса"
// @Success 200 {object} TreeResponse "Изменённое дерево"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 401 {object} corehttp.ErrorResponse "Unauthorized"
// @Failure 403 {object} corehttp.ErrorResponse "Forbidden"
// @Failure 404 {object} corehttp.ErrorResponse "Tree not found"
// @Router /trees/{tree_id} [patch]
func (h *TreesHTTPHandler) PatchTree(w http.ResponseWriter, r *http.Request) {
	userID, err := request.GetUserID(r)
	if err != nil {
		corehttp.Error(w, err, "authentication is required")
		return
	}
	treeID, err := request.GetTreeID(r)
	if err != nil {
		corehttp.Error(w, err, "tree id is invalid")
		return
	}
	var input PatchTreeRequest
	if err := request.DecodeJSON(r, &input); err != nil {
		corehttp.Error(w, err, "request body contains invalid JSON")
		return
	}
	tree, err := h.treesService.PatchTree(r.Context(), userID, treeID, domain.PatchTreeInput{Name: input.Name})
	if err != nil {
		corehttp.Error(w, err, "could not update tree")
		return
	}
	corehttp.JSON(w, http.StatusOK, tree)
}
