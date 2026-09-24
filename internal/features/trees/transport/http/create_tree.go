package http

import (
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// CreateTree godoc
// @Summary Создание дерева
// @Description Создание дерева с автоматическим назначением текущего пользователя владельцем
// @Tags trees
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTreeRequest true "CreateTree тело запроса"
// @Success 201 {object} TreeResponse "Успешно созданное дерево"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 401 {object} corehttp.ErrorResponse "Unauthorized"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /trees [post]
func (h *TreesHTTPHandler) CreateTree(w http.ResponseWriter, r *http.Request) {
	userID, err := request.GetUserID(r)
	if err != nil {
		corehttp.Error(w, err, "authentication is required")
		return
	}
	var input CreateTreeRequest
	if err := request.DecodeJSON(r, &input); err != nil {
		corehttp.Error(w, err, "request body contains invalid JSON")
		return
	}
	tree, err := h.treesService.CreateTree(r.Context(), userID, domain.CreateTreeInput{Name: input.Name})
	if err != nil {
		corehttp.Error(w, err, "could not create tree")
		return
	}
	corehttp.JSON(w, http.StatusCreated, tree)
}
