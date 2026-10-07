package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// GetFamilies godoc
// @Summary Список семей
// @Tags families
// @Security BearerAuth
// @Produce json
// @Param tree_id path string true "ULID дерева"
// @Success 200 {array} FamilyResponse
// @Router /trees/{tree_id}/families [get]
func (h *FamiliesHTTPHandler) GetFamilies(w http.ResponseWriter, r *http.Request) {
	userID, treeID, err := request.GetTreeContext(r)
	if err != nil {
		corehttp.Error(w, err, "tree access is invalid")
		return
	}
	families, err := h.familiesService.GetFamilies(r.Context(), userID, treeID)
	if err != nil {
		corehttp.Error(w, err, "could not get families")
		return
	}
	corehttp.JSON(w, http.StatusOK, families)
}
