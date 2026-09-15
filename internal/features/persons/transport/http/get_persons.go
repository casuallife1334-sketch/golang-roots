package http

import (
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// GetPersons godoc
// @Summary Список людей
// @Description Получение списка людей в генеалогическом дереве
// @Tags persons
// @Produce json
// @Success 200 {array} PersonResponse "Успешное получение списка людей"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /persons [get]
func (h *PersonsHTTPHandler) GetPersons(w http.ResponseWriter, r *http.Request) {
	items, err := h.personsService.GetPersons(r.Context())
	if err != nil {
		corehttp.Error(w, err, "could not get persons")
		return
	}
	corehttp.JSON(w, http.StatusOK, items)
}
