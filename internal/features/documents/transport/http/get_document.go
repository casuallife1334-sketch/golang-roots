package http

import (
	"errors"
	"fmt"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/documents/repository"
	"io"
	"mime"
	"net/http"
)

// GetDocument godoc
// @Summary Скачивание документа
// @Description Получение бинарного содержимого документа
// @Tags documents
// @Security BearerAuth
// @Produce application/octet-stream
// @Param tree_id path string true "ULID дерева"
// @Param id path string true "ULID документа"
// @Success 200 {file} binary "Документ"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 404 {object} corehttp.ErrorResponse "Document not found"
// @Router /trees/{tree_id}/documents/{id} [get]
func (h *DocumentsHTTPHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
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
	document, file, err := h.documentsService.OpenDocument(r.Context(), userID, treeID, id)
	if errors.Is(err, repository.ErrNotFound) {
		corehttp.Error(w, err, "the requested document was not found")
		return
	}
	if err != nil {
		corehttp.Error(w, err, "could not get document")
		return
	}
	defer file.Close()
	w.Header().Set("Content-Type", document.ContentType)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": document.FileName}))
	w.Header().Set("Content-Length", fmt.Sprint(document.SizeBytes))
	_, _ = io.Copy(w, file)
}
