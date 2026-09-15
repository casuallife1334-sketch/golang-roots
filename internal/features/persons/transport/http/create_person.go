package http

import (
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// CreatePerson godoc
// @Summary 	Создание человека
// @Description Создание нового человека в генеалогическом дереве
// @Tags 		persons
// @Accept 		json
// @Produce 	json
// @Param 		request body CreatePersonRequest true "CreatePerson тело запроса"
// @Success 	201 {object} PersonResponse "Успешно созданный человек"
// @Failure 	400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 	500 {object} corehttp.ErrorResponse "internal server error"
// @Router 		/persons [post]
func (h *PersonsHTTPHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var input domain.CreatePersonInput
	if err := request.DecodeJSON(r, &input); err != nil {
		corehttp.Error(w, err, "request body contains invalid JSON")
		return
	}
	p, err := h.personsService.CreatePerson(r.Context(), input)
	if err != nil {
		corehttp.Error(w, err, "could not create person")
		return
	}
	corehttp.JSON(w, http.StatusCreated, p)
}
