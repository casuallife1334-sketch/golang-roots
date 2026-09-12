package http

import (
	"errors"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/persons/repository"
	"genealogy-tree/internal/features/persons/service"
	"net/http"
)

func (h *PersonsHTTPHandler) PatchPerson(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetULIDPathValue(r, "id")
	if err != nil {
		corehttp.Error(w, err, "person id is invalid")
		return
	}
	var input domain.PatchPersonInput
	if err := request.DecodeJSON(r, &input); err != nil {
		corehttp.Error(w, err, "request body contains invalid JSON")
		return
	}
	p, err := h.personsService.PatchPerson(r.Context(), id, input)
	if errors.Is(err, repository.ErrNotFound) {
		corehttp.Error(w, err, "the requested person was not found")
		return
	}
	if errors.Is(err, service.ErrInvalid) {
		corehttp.Error(w, err, "person data is invalid")
		return
	}
	if err != nil {
		corehttp.Error(w, err, "could not update person")
		return
	}
	corehttp.JSON(w, http.StatusOK, p)
}
