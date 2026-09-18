package http

import (
	"errors"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/persons/repository"
	"net/http"
)

// GetPerson godoc
// @Summary Получение человека
// @Description Получение конкретного человека по ULID
// @Tags persons
// @Security BearerAuth
// @Produce json
// @Param tree_id path string true "ULID дерева"
// @Param id path string true "ULID получаемого человека"
// @Success 200 {object} PersonResponse "Человек успешно найден"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 404 {object} corehttp.ErrorResponse "Person not found"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /trees/{tree_id}/persons/{id} [get]
func (h *PersonsHTTPHandler) GetPerson(w http.ResponseWriter, r *http.Request) {
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
	p, err := h.personsService.GetPerson(r.Context(), userID, treeID, id)
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
