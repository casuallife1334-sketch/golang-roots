package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

const maxPhotoSize = 10 << 20

// UploadPersonPhoto godoc
// @Summary Загрузка фотографии человека
// @Description Загрузка или замена фотографии человека, максимальный размер 10 MB
// @Tags persons
// @Security BearerAuth
// @Accept mpfd
// @Produce json
// @Param tree_id path string true "ULID дерева"
// @Param id path string true "ULID человека"
// @Param file formData file true "Фотография человека"
// @Success 200 {object} PersonResponse "Фотография успешно загружена"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 404 {object} corehttp.ErrorResponse "Person not found"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /trees/{tree_id}/persons/{id}/photo [post]
func (h *PersonsHTTPHandler) UploadPersonPhoto(w http.ResponseWriter, r *http.Request) {
	userID, treeID, err := getTreeContext(r)
	if err != nil {
		corehttp.Error(w, err, "tree access is invalid")
		return
	}
	id, err := request.GetULIDPathValue(r, "id")
	if err != nil {
		corehttp.Error(w, err, "person id is invalid")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxPhotoSize)
	if err := r.ParseMultipartForm(maxPhotoSize); err != nil {
		corehttp.Error(w, err, "photo must be a multipart file up to 10 MB")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		corehttp.Error(w, err, "multipart field 'file' is required")
		return
	}
	defer file.Close()
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	person, err := h.personsService.UploadPersonPhoto(r.Context(), userID, treeID, id, contentType, file)
	if err != nil {
		corehttp.Error(w, err, "could not upload person photo")
		return
	}
	corehttp.JSON(w, http.StatusOK, person)
}
