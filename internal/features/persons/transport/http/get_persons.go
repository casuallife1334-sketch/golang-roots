package http

import (
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

func (h *PersonsHTTPHandler) GetPersons(w http.ResponseWriter, r *http.Request) {
	items, err := h.personsService.GetPersons(r.Context())
	if err != nil {
		corehttp.Error(w, err, "could not get persons")
		return
	}
	corehttp.JSON(w, http.StatusOK, items)
}
