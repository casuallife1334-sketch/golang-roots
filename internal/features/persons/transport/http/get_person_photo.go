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

// GetPersonPhoto godoc
// @Summary Получение фотографии человека
// @Description Получение бинарного содержимого фотографии человека
// @Tags persons
// @Security BearerAuth
// @Produce application/octet-stream
// @Param tree_id path string true "ULID дерева"
// @Param id path string true "ULID человека"
// @Success 200 {file} binary "Фотография человека"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 404 {object} corehttp.ErrorResponse "Person photo not found"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /trees/{tree_id}/persons/{id}/photo [get]
func (h *PersonsHTTPHandler) GetPersonPhoto(w http.ResponseWriter, r *http.Request) {
	userID, treeID, err := getTreeContext(r)
	if err != nil {
		corehttp.Error(w, err, "tree access is invalid")
		return
	}
	id, err := request.GetULIDPathValue(r, "id")
	if err != nil {
		corehttp.Error(w, err, "person id is invalid")
		return
	}
	photo, err := h.personsService.GetPersonPhoto(r.Context(), userID, treeID, id)
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
