package http

import (
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// CreateFamily godoc
// @Summary Создание семьи
// @Description Создание пустого контейнера семьи
// @Tags families
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param tree_id path string true "ULID дерева"
// @Param request body CreateFamilyRequest true "CreateFamily тело запроса"
// @Success 201 {object} FamilyResponse
// @Router /trees/{tree_id}/families [post]
func (h *FamiliesHTTPHandler) CreateFamily(w http.ResponseWriter, r *http.Request) {
	userID, treeID, err := request.GetTreeContext(r)
	if err != nil {
		corehttp.Error(w, err, "tree access is invalid")
		return
	}
	var body CreateFamilyRequest
	if err := request.DecodeJSON(r, &body); err != nil {
		corehttp.Error(w, err, "request body contains invalid JSON")
		return
	}
	family, err := h.familiesService.CreateFamily(r.Context(), userID, treeID, domain.CreateFamilyInput{
		Name: body.Name, Metadata: body.Metadata,
	})
	if err != nil {
		corehttp.Error(w, err, "could not create family")
		return
	}
	corehttp.JSON(w, http.StatusCreated, family)
}
