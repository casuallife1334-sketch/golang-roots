package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

func (h *FamiliesHTTPHandler) AttachRelationship(w http.ResponseWriter, r *http.Request) {
	userID, treeID, familyID, err := familyContext(r)
	if err != nil {
		corehttp.Error(w, err, "family context is invalid")
		return
	}
	var body AttachRelationshipRequest
	if err := request.DecodeJSON(r, &body); err != nil {
		corehttp.Error(w, err, "request body contains invalid JSON")
		return
	}
	if err := h.familiesService.AttachRelationship(r.Context(), userID, treeID, familyID, body.RelationshipID); err != nil {
		corehttp.Error(w, err, "could not attach relationship")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
