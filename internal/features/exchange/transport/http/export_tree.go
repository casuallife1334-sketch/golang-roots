package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	"genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// ExportTree godoc
// @Summary Экспорт дерева в GEDCOM
// @Description Экспортирует людей и родственные связи выбранного дерева в GEDCOM 5.5.1
// @Tags exchange
// @Security BearerAuth
// @Produce application/x-gedcom
// @Param tree_id path string true "ULID дерева"
// @Param format query string true "Формат файла" Enums(gedcom)
// @Success 200 {file} binary "GEDCOM-файл"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Router /trees/{tree_id}/export [get]
func (h *ExchangeHTTPHandler) ExportTree(w http.ResponseWriter, r *http.Request) {
	if err := requireGEDCOMFormat(r); err != nil {
		response.Error(w, err, "format must be gedcom")
		return
	}
	userID, treeID, err := request.GetTreeContext(r)
	if err != nil {
		response.Error(w, err, "tree access is invalid")
		return
	}
	content, err := h.exchangeService.Export(r.Context(), userID, treeID)
	if err != nil {
		response.Error(w, err, "could not export tree")
		return
	}
	w.Header().Set("Content-Type", "application/x-gedcom; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="tree.ged"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}
