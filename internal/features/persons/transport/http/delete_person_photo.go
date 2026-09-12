package http

import (
	"errors"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/persons/repository"
	"net/http"
)

func (h *PersonsHTTPHandler) DeletePersonPhoto(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetULIDPathValue(r, "id")
	if err != nil {
		corehttp.Error(w, err, "person id is invalid")
		return
	}
	err = h.personsService.DeletePersonPhoto(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		corehttp.Error(w, err, "person was not found")
		return
	}
	if err != nil {
		corehttp.Error(w, err, "could not delete person photo")
		return
	}
	corehttp.NoContent(w)
}
