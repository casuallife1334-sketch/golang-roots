package http

import (
	"errors"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/relationships/repository"
	"genealogy-tree/internal/features/relationships/service"
	"net/http"
)

// PatchRelationship godoc
// @Summary Изменение связи
// @Description Изменение метаданных существующей связи
// @Tags relationships
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param tree_id path string true "ULID дерева"
// @Param id path string true "ULID изменяемой связи"
// @Param request body PatchRelationshipRequest true "PatchRelationship тело запроса"
// @Success 200 {object} RelationshipResponse "Успешно изменённая связь"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 404 {object} corehttp.ErrorResponse "Relationship not found"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /trees/{tree_id}/relationships/{id} [patch]
func (h *RelationshipsHTTPHandler) PatchRelationship(w http.ResponseWriter, r *http.Request) {
	userID, treeID, err := request.GetTreeContext(r)
	if err != nil {
		corehttp.Error(w, err, "tree access is invalid")
		return
	}
	id, err := request.GetULIDPathValue(r, "id")
	if err != nil {
		corehttp.Error(w, err, "relationship id is invalid")
		return
	}
	var requestBody PatchRelationshipRequest
	if err := request.DecodeJSON(r, &requestBody); err != nil {
		corehttp.Error(w, err, "request body contains invalid JSON")
		return
	}
	rel, err := h.relationshipsService.PatchRelationship(r.Context(), userID, treeID, id, domain.PatchRelationshipInput{Metadata: requestBody.Metadata})
	if errors.Is(err, repository.ErrNotFound) {
		corehttp.Error(w, err, "the requested relationship was not found")
		return
	}
	if errors.Is(err, service.ErrInvalid) {
		corehttp.Error(w, err, "relationship data is invalid")
		return
	}
	if err != nil {
		corehttp.Error(w, err, "could not update relationship")
		return
	}
	corehttp.JSON(w, http.StatusOK, rel)
}
