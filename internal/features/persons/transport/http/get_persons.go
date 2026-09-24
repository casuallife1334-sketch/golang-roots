package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// GetPersons godoc
// @Summary Список людей
// @Description Получение списка людей в генеалогическом дереве
// @Tags persons
// @Security BearerAuth
// @Produce json
// @Param tree_id path string true "ULID дерева"
// @Success 200 {array} PersonResponse "Успешное получение списка людей"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /trees/{tree_id}/persons [get]
func (h *PersonsHTTPHandler) GetPersons(w http.ResponseWriter, r *http.Request) {
	userID, treeID, err := request.GetTreeContext(r)
	if err != nil {
		corehttp.Error(w, err, "tree access is invalid")
		return
	}
	items, err := h.personsService.GetPersons(r.Context(), userID, treeID)
	if err != nil {
		corehttp.Error(w, err, "could not get persons")
		return
	}
	corehttp.JSON(w, http.StatusOK, items)
}
