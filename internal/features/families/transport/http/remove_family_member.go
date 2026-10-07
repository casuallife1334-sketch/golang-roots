package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

func (h *FamiliesHTTPHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	userID, treeID, familyID, err := familyContext(r)
	if err != nil {
		corehttp.Error(w, err, "family context is invalid")
		return
	}
	personID, err := request.GetULIDPathValue(r, "person_id")
	if err != nil {
		corehttp.Error(w, err, "person id is invalid")
		return
	}
	if err := h.familiesService.RemoveMember(r.Context(), userID, treeID, familyID, personID); err != nil {
		corehttp.Error(w, err, "could not remove family member")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
