package http

import (
	"errors"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/persons/repository"
	"net/http"
)

// DeletePerson godoc
// @Summary Удаление человека
// @Description Удаление человека и связанных с ним relationships
// @Tags persons
// @Security BearerAuth
// @Param tree_id path string true "ULID дерева"
// @Param id path string true "ULID удаляемого человека"
// @Success 204 "Успешное удаление человека"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 404 {object} corehttp.ErrorResponse "Person not found"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /trees/{tree_id}/persons/{id} [delete]
func (h *PersonsHTTPHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
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
	err = h.personsService.DeletePerson(r.Context(), userID, treeID, id)
	if errors.Is(err, repository.ErrNotFound) {
		corehttp.Error(w, err, "the requested person was not found")
		return
	}
	if err != nil {
		corehttp.Error(w, err, "could not delete person")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
