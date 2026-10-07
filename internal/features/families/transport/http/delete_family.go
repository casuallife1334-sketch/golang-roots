package http

import (
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// DeleteFamily godoc
// @Summary Удаление семьи
// @Tags families
// @Security BearerAuth
// @Param tree_id path string true "ULID дерева"
// @Param family_id path string true "ULID семьи"
// @Success 204
// @Router /trees/{tree_id}/families/{family_id} [delete]
func (h *FamiliesHTTPHandler) DeleteFamily(w http.ResponseWriter, r *http.Request) {
	userID, treeID, familyID, err := familyContext(r)
	if err != nil {
		corehttp.Error(w, err, "family context is invalid")
		return
	}
	if err := h.familiesService.DeleteFamily(r.Context(), userID, treeID, familyID); err != nil {
		corehttp.Error(w, err, "could not delete family")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
