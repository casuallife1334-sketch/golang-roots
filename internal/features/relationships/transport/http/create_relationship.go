package http

import (
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

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
