package http

import (
	"errors"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/relationships/repository"
	"net/http"
)

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
