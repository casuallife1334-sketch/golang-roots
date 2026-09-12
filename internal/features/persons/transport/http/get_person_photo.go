package http

import (
	"errors"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/persons/repository"
	"genealogy-tree/internal/features/persons/service"
	"io"
	"net/http"
)

func (h *PersonsHTTPHandler) GetPersonPhoto(w http.ResponseWriter, r *http.Request) {
	id, err := request.GetULIDPathValue(r, "id")
	if err != nil {
		corehttp.Error(w, err, "person id is invalid")
		return
	}
	photo, err := h.personsService.GetPersonPhoto(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) || errors.Is(err, service.ErrPhotoNotFound) {
		corehttp.Error(w, err, "person photo was not found")
		return
	}
	if err != nil {
		corehttp.Error(w, err, "could not get person photo")
		return
	}
	defer photo.Close()
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := io.Copy(w, photo); err != nil {
		return
	}
}
