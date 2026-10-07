package http

import (
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

func (h *FamiliesHTTPHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	userID, treeID, familyID, err := familyContext(r)
	if err != nil {
		corehttp.Error(w, err, "family context is invalid")
		return
	}
	var body AddFamilyMemberRequest
	if err := request.DecodeJSON(r, &body); err != nil {
		corehttp.Error(w, err, "request body contains invalid JSON")
		return
	}
	err = h.familiesService.AddMember(r.Context(), userID, treeID, familyID, domain.AddFamilyMemberInput{
		PersonID: body.PersonID, Role: body.Role,
	})
	if err != nil {
		corehttp.Error(w, err, "could not add family member")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
