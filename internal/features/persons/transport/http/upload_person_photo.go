package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

const maxPhotoSize = 10 << 20

func (h *PersonsHTTPHandler) UploadPersonPhoto(w http.ResponseWriter, r *http.Request) {
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
	person, err := h.personsService.UploadPersonPhoto(r.Context(), id, contentType, file)
	if err != nil {
		corehttp.Error(w, err, "could not upload person photo")
		return
	}
	corehttp.JSON(w, http.StatusOK, person)
}
