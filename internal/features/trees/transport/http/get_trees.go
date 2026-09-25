package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// GetTrees godoc
// @Summary Список деревьев
// @Description Получение деревьев, доступных текущему пользователю
// @Tags trees
// @Produce json
// @Security BearerAuth
// @Success 200 {array} TreeResponse "Список деревьев"
// @Failure 401 {object} corehttp.ErrorResponse "Unauthorized"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /trees [get]
func (h *TreesHTTPHandler) GetTrees(w http.ResponseWriter, r *http.Request) {
	userID, err := request.GetUserID(r)
	if err != nil {
		corehttp.Error(w, err, "authentication is required")
		return
	}
	trees, err := h.treesService.GetTrees(r.Context(), userID)
	if err != nil {
		corehttp.Error(w, err, "could not get trees")
		return
	}
	corehttp.JSON(w, http.StatusOK, trees)
}
