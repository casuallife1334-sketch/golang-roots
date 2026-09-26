package http

import (
	"errors"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/documents/repository"
	"net/http"
)

// DeleteDocument godoc
// @Summary Удаление документа
// @Description Удаление документа из дерева и файлового хранилища
// @Tags documents
// @Security BearerAuth
// @Param tree_id path string true "ULID дерева"
// @Param id path string true "ULID документа"
// @Success 204 "Документ удалён"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 404 {object} corehttp.ErrorResponse "Document not found"
// @Router /trees/{tree_id}/documents/{id} [delete]
func (h *DocumentsHTTPHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	userID, treeID, err := request.GetTreeContext(r)
	if err != nil {
		corehttp.Error(w, err, "tree access is invalid")
		return
	}
	id, err := request.GetULIDPathValue(r, "id")
	if err != nil {
		corehttp.Error(w, err, "document id is invalid")
		return
	}
	if err := h.documentsService.DeleteDocument(r.Context(), userID, treeID, id); errors.Is(err, repository.ErrNotFound) {
		corehttp.Error(w, err, "the requested document was not found")
		return
	} else if err != nil {
		corehttp.Error(w, err, "could not delete document")
		return
	}
	corehttp.NoContent(w)
}
