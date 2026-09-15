package http

import (
	"errors"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/relationships/repository"
	"net/http"
)

// GetRelationship godoc
// @Summary Получение связи
// @Description Получение конкретной связи по ULID
// @Tags relationships
// @Produce json
// @Param id path string true "ULID получаемой связи"
// @Success 200 {object} RelationshipResponse "Связь успешно найдена"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 404 {object} corehttp.ErrorResponse "Relationship not found"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /relationships/{id} [get]
func (h *RelationshipsHTTPHandler) GetRelationship(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetULIDPathValue(r, "id")
	if err != nil {
		corehttp.Error(w, err, "relationship id is invalid")
		return
	}
	rel, err := h.relationshipsService.GetRelationship(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		corehttp.Error(w, err, "the requested relationship was not found")
		return
	}
	if err != nil {
		corehttp.Error(w, err, "could not get relationship")
		return
	}
	corehttp.JSON(w, http.StatusOK, rel)
}
