package http

import (
	"errors"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/trees/repository"
	"net/http"
)

// GetTree godoc
// @Summary Получение дерева
// @Description Получение доступного пользователю дерева
// @Tags trees
// @Produce json
// @Security BearerAuth
// @Param tree_id path string true "ULID дерева"
// @Success 200 {object} TreeResponse "Дерево"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 401 {object} corehttp.ErrorResponse "Unauthorized"
// @Failure 404 {object} corehttp.ErrorResponse "Tree not found"
// @Router /trees/{tree_id} [get]
func (h *TreesHTTPHandler) GetTree(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		corehttp.Error(w, err, "authentication is required")
		return
	}
	treeID, err := getTreeID(r)
	if err != nil {
		corehttp.Error(w, err, "tree id is invalid")
		return
	}
	tree, err := h.treesService.GetTree(r.Context(), userID, treeID)
	if errors.Is(err, repository.ErrNotFound) {
		corehttp.Error(w, err, "the requested tree was not found")
		return
	}
	if err != nil {
		corehttp.Error(w, err, "could not get tree")
		return
	}
	corehttp.JSON(w, http.StatusOK, tree)
}
