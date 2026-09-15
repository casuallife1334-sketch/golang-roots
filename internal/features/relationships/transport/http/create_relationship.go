package http

import (
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// CreateRelationship godoc
// @Summary Создание связи
// @Description Создание новой связи между двумя людьми
// @Tags relationships
// @Accept json
// @Produce json
// @Param request body CreateRelationshipRequest true "CreateRelationship тело запроса"
// @Success 201 {object} RelationshipResponse "Успешно созданная связь"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 404 {object} corehttp.ErrorResponse "Person not found"
// @Failure 409 {object} corehttp.ErrorResponse "Conflict"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /relationships [post]
func (h *RelationshipsHTTPHandler) CreateRelationship(w http.ResponseWriter, r *http.Request) {
	var in domain.CreateRelationshipInput
	if err := request.DecodeJSON(r, &in); err != nil {
		corehttp.Error(w, err, "request body contains invalid JSON")
		return
	}
	rel, err := h.relationshipsService.CreateRelationship(r.Context(), in)
	if err != nil {
		corehttp.Error(w, err, "could not create relationship")
		return
	}
	corehttp.JSON(w, http.StatusCreated, rel)
}
