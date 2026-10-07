package http

import (
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// GetFamily godoc
// @Summary Получение семьи
// @Tags families
// @Security BearerAuth
// @Produce json
// @Param tree_id path string true "ULID дерева"
// @Param family_id path string true "ULID семьи"
// @Success 200 {object} FamilyResponse
// @Router /trees/{tree_id}/families/{family_id} [get]
func (h *FamiliesHTTPHandler) GetFamily(w http.ResponseWriter, r *http.Request) {
	userID, treeID, familyID, err := familyContext(r)
	if err != nil {
		corehttp.Error(w, err, "family context is invalid")
		return
	}
	family, err := h.familiesService.GetFamily(r.Context(), userID, treeID, familyID)
	if err != nil {
		corehttp.Error(w, err, "could not get family")
		return
	}
	corehttp.JSON(w, http.StatusOK, family)
}
