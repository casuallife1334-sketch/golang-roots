package http

import (
	"errors"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/relationships/repository"
	"net/http"
)

// DeleteRelationship godoc
// @Summary Удаление связи
// @Description Удаление существующей связи по ULID
// @Tags relationships
// @Param id path string true "ULID удаляемой связи"
// @Success 204 "Успешное удаление связи"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 404 {object} corehttp.ErrorResponse "Relationship not found"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /relationships/{id} [delete]
func (h *RelationshipsHTTPHandler) DeleteRelationship(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetULIDPathValue(r, "id")
	if err != nil {
		corehttp.Error(w, err, "relationship id is invalid")
		return
	}
	err = h.relationshipsService.DeleteRelationship(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		corehttp.Error(w, err, "the requested relationship was not found")
		return
	}
	if err != nil {
		corehttp.Error(w, err, "could not delete relationship")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
