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

// PatchPerson godoc
// @Summary Изменение человека
// @Description Изменение информации о существующем человеке
// @Tags persons
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param tree_id path string true "ULID дерева"
// @Param id path string true "ULID изменяемого человека"
// @Param request body domain.PatchPersonInput true "PatchPerson тело запроса"
// @Success 200 {object} PersonResponse "Успешно изменённый человек"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 404 {object} corehttp.ErrorResponse "Person not found"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /trees/{tree_id}/persons/{id} [patch]
func (h *PersonsHTTPHandler) PatchPerson(w http.ResponseWriter, r *http.Request) {
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
	var input domain.PatchPersonInput
	if err := request.DecodeJSON(r, &input); err != nil {
		corehttp.Error(w, err, "request body contains invalid JSON")
		return
	}
	p, err := h.personsService.PatchPerson(r.Context(), userID, treeID, id, input)
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
