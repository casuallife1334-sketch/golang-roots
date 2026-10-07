package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

func (h *FamiliesHTTPHandler) DetachRelationship(w http.ResponseWriter, r *http.Request) {
	userID, treeID, familyID, err := familyContext(r)
	if err != nil {
		corehttp.Error(w, err, "family context is invalid")
		return
	}
	relationshipID, err := request.GetULIDPathValue(r, "relationship_id")
	if err != nil {
		corehttp.Error(w, err, "relationship id is invalid")
		return
	}
	if err := h.familiesService.DetachRelationship(r.Context(), userID, treeID, familyID, relationshipID); err != nil {
		corehttp.Error(w, err, "could not detach relationship")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
