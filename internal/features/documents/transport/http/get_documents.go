package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// GetDocuments godoc
// @Summary Список документов
// @Description Получение документов человека или связи
// @Tags documents
// @Security BearerAuth
// @Produce json
// @Param tree_id path string true "ULID дерева"
// @Param owner_type query string true "Тип владельца: person или relationship"
// @Param owner_id query string true "ULID владельца документа"
// @Success 200 {array} DocumentResponse "Список документов"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Router /trees/{tree_id}/documents [get]
func (h *DocumentsHTTPHandler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	userID, treeID, err := request.GetTreeContext(r)
	if err != nil {
		corehttp.Error(w, err, "tree access is invalid")
		return
	}
	owner, err := documentOwner(r)
	if err != nil {
		corehttp.Error(w, err, "document owner is invalid")
		return
	}
	documents, err := h.documentsService.GetDocuments(r.Context(), userID, treeID, owner)
	if err != nil {
		corehttp.Error(w, err, "could not get documents")
		return
	}
	corehttp.JSON(w, http.StatusOK, documents)
}
