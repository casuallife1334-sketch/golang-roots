package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// GetRelationships godoc
// @Summary Список связей
// @Description Получение всех связей с опциональной фильтрацией по ULID человека
// @Tags relationships
// @Security BearerAuth
// @Produce json
// @Param tree_id path string true "ULID дерева"
// @Param person_id query string false "Фильтрация связей по ULID человека"
// @Success 200 {array} RelationshipResponse "Успешное получение списка связей"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /trees/{tree_id}/relationships [get]
func (h *RelationshipsHTTPHandler) GetRelationships(w http.ResponseWriter, r *http.Request) {
	userID, treeID, err := getTreeContext(r)
	if err != nil {
		corehttp.Error(w, err, "tree access is invalid")
		return
	}
	personID, err := request.GetOptionalULIDQueryValue(r, "person_id")
	if err != nil {
		corehttp.Error(w, err, "person_id is invalid")
		return
	}
	relationships, err := h.relationshipsService.GetRelationships(r.Context(), userID, treeID, personID)
	if err != nil {
		corehttp.Error(w, err, "could not get relationships")
		return
	}
	corehttp.JSON(w, http.StatusOK, relationships)
}
