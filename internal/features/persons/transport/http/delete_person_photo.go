package http

import (
	"errors"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/persons/repository"
	"net/http"
)

// DeletePersonPhoto godoc
// @Summary Удаление фотографии человека
// @Description Удаление фотографии человека из хранилища
// @Tags persons
// @Security BearerAuth
// @Param tree_id path string true "ULID дерева"
// @Param id path string true "ULID человека"
// @Success 204 "Фотография успешно удалена"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 404 {object} corehttp.ErrorResponse "Person not found"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /trees/{tree_id}/persons/{id}/photo [delete]
func (h *PersonsHTTPHandler) DeletePersonPhoto(w http.ResponseWriter, r *http.Request) {
	userID, treeID, err := request.GetTreeContext(r)
	if err != nil {
		corehttp.Error(w, err, "tree access is invalid")
		return
	}
	id, err := request.GetULIDPathValue(r, "id")
	if err != nil {
		corehttp.Error(w, err, "person id is invalid")
		return
	}
	err = h.personsService.DeletePersonPhoto(r.Context(), userID, treeID, id)
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
