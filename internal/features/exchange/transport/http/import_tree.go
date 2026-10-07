package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	"genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// ImportTree godoc
// @Summary Импорт GEDCOM в дерево
// @Description Проверяет и атомарно импортирует людей и родственные связи в выбранное дерево
// @Tags exchange
// @Security BearerAuth
// @Accept mpfd
// @Produce json
// @Param tree_id path string true "ULID дерева"
// @Param format query string true "Формат файла" Enums(gedcom)
// @Param file formData file true "GEDCOM-файл"
// @Success 201 {object} ImportResponse "Результат импорта"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 413 {object} response.ErrorResponse "Payload Too Large"
// @Router /trees/{tree_id}/import [post]
func (h *ExchangeHTTPHandler) ImportTree(w http.ResponseWriter, r *http.Request) {
	if err := requireGEDCOMFormat(r); err != nil {
		response.Error(w, err, "format must be gedcom")
		return
	}
	userID, treeID, err := request.GetTreeContext(r)
	if err != nil {
		response.Error(w, err, "tree access is invalid")
		return
	}
	file, err := readGEDCOMFile(w, r)
	if err != nil {
		response.Error(w, err, "could not read GEDCOM file")
		return
	}
	defer file.Close()
	result, _, err := h.exchangeService.Import(r.Context(), userID, treeID, file)
	if err != nil {
		response.Error(w, err, "could not import GEDCOM file")
		return
	}
	response.JSON(w, http.StatusCreated, ImportResponse{
		Format:        "gedcom",
		Persons:       result.Persons,
		Relationships: result.Relationships,
	})
}
