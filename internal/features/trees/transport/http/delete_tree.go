package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// DeleteTree godoc
// @Summary Удаление дерева
// @Description Удаление дерева владельцем вместе с его участниками и связанными данными
// @Tags trees
// @Security BearerAuth
// @Param tree_id path string true "ULID дерева"
// @Success 204 "Дерево удалено"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 401 {object} corehttp.ErrorResponse "Unauthorized"
// @Failure 403 {object} corehttp.ErrorResponse "Forbidden"
// @Failure 404 {object} corehttp.ErrorResponse "Tree not found"
// @Router /trees/{tree_id} [delete]
func (h *TreesHTTPHandler) DeleteTree(w http.ResponseWriter, r *http.Request) {
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
	if err := h.treesService.DeleteTree(r.Context(), userID, treeID); err != nil {
		corehttp.Error(w, err, "could not delete tree")
		return
	}
	corehttp.NoContent(w)
}
