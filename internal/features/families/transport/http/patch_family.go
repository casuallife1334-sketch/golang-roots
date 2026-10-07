package http

import (
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// PatchFamily godoc
// @Summary Изменение семьи
// @Tags families
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param tree_id path string true "ULID дерева"
// @Param family_id path string true "ULID семьи"
// @Param request body PatchFamilyRequest true "PatchFamily тело запроса"
// @Success 200 {object} FamilyResponse
// @Router /trees/{tree_id}/families/{family_id} [patch]
func (h *FamiliesHTTPHandler) PatchFamily(w http.ResponseWriter, r *http.Request) {
	userID, treeID, familyID, err := familyContext(r)
	if err != nil {
		corehttp.Error(w, err, "family context is invalid")
		return
	}
	var body PatchFamilyRequest
	if err := request.DecodeJSON(r, &body); err != nil {
		corehttp.Error(w, err, "request body contains invalid JSON")
		return
	}
	var metadata *map[string]any
	if body.Metadata.Set {
		metadata = body.Metadata.Value
	}
	family, err := h.familiesService.PatchFamily(r.Context(), userID, treeID, familyID, domain.PatchFamilyInput{
		Name: body.Name, Metadata: metadata,
	})
	if err != nil {
		corehttp.Error(w, err, "could not update family")
		return
	}
	corehttp.JSON(w, http.StatusOK, family)
}
