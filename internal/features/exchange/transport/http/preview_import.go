package http

import (
	"genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// PreviewImport godoc
// @Summary Предварительный просмотр GEDCOM
// @Description Проверяет GEDCOM-файл без записи в выбранное дерево
// @Tags exchange
// @Security BearerAuth
// @Accept mpfd
// @Produce json
// @Param format query string true "Формат файла" Enums(gedcom)
// @Param file formData file true "GEDCOM-файл"
// @Success 200 {object} PreviewResponse "Результат проверки"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 413 {object} response.ErrorResponse "Payload Too Large"
// @Router /imports/preview [post]
func (h *ExchangeHTTPHandler) PreviewImport(w http.ResponseWriter, r *http.Request) {
	if err := requireGEDCOMFormat(r); err != nil {
		response.Error(w, err, "format must be gedcom")
		return
	}
	file, err := readGEDCOMFile(w, r)
	if err != nil {
		response.Error(w, err, "could not read GEDCOM file")
		return
	}
	defer file.Close()
	snapshot, err := h.exchangeService.Preview(file)
	if err != nil {
		response.Error(w, err, "could not parse GEDCOM file")
		return
	}
	response.JSON(w, http.StatusOK, PreviewResponse{
		Format:        "gedcom",
		Persons:       len(snapshot.Persons),
		Relationships: len(snapshot.Relationships),
		Warnings:      snapshot.Warnings,
		Errors:        snapshot.Errors,
	})
}
