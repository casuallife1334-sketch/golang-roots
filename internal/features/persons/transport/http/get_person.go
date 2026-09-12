package http

import (
	"errors"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/persons/repository"
	"net/http"
)

func (h *PersonsHTTPHandler) GetPerson(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetULIDPathValue(r, "id")
	if err != nil {
		corehttp.Error(w, err, "person id is invalid")
		return
	}
	p, err := h.personsService.GetPerson(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		corehttp.Error(w, err, "the requested person was not found")
		return
	}
	if err != nil {
		corehttp.Error(w, err, "could not get person")
		return
	}
	corehttp.JSON(w, http.StatusOK, p)
}
