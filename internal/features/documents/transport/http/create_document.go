package http

import (
	"errors"
	"fmt"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/core/validation"
	"genealogy-tree/internal/features/documents/service"
	"io"
	"mime/multipart"
	"net/http"
)

const maxDocumentSize = validation.MaxDocumentSize

// CreateDocument godoc
// @Summary Загрузка документа
// @Description Загрузка документа к человеку или связи
// @Tags documents
// @Security BearerAuth
// @Accept mpfd
// @Produce json
// @Param tree_id path string true "ULID дерева"
// @Param owner_type query string true "Тип владельца: person или relationship"
// @Param owner_id query string true "ULID владельца документа"
// @Param file formData file true "Документ"
// @Success 201 {object} DocumentResponse "Документ загружен"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 403 {object} corehttp.ErrorResponse "Forbidden"
// @Failure 413 {object} corehttp.ErrorResponse "Payload Too Large"
// @Failure 404 {object} corehttp.ErrorResponse "Owner not found"
// @Router /trees/{tree_id}/documents [post]
func (h *DocumentsHTTPHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
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
	r.Body = http.MaxBytesReader(w, r.Body, maxDocumentSize+(1<<20))
	if err := r.ParseMultipartForm(maxDocumentSize + (1 << 20)); err != nil {
		corehttp.Error(w, err, "document must be a multipart file up to 20 MB")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		corehttp.Error(w, err, "multipart field 'file' is required")
		return
	}
	defer file.Close()
	if header.Size > maxDocumentSize {
		corehttp.Error(w, service.ErrTooLarge, "document must be no larger than 20 MB")
		return
	}
	contentType, err := detectContentType(file, header.Header.Get("Content-Type"))
	if err != nil {
		corehttp.Error(w, service.ErrInvalid, "document type is not supported")
		return
	}
	document, err := h.documentsService.CreateDocument(r.Context(), userID, treeID, service.UploadInput{
		Owner:       owner,
		FileName:    header.Filename,
		ContentType: contentType,
		SizeBytes:   header.Size,
		File:        file,
	})
	if err != nil {
		corehttp.Error(w, err, "could not upload document")
		return
	}
	corehttp.JSON(w, http.StatusCreated, document)
}

func detectContentType(file multipart.File, provided string) (string, error) {
	buffer := make([]byte, 512)
	read, err := file.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	detected := http.DetectContentType(buffer[:read])
	if validation.IsDocumentContentType(detected) {
		return detected, nil
	}
	if detected == "application/zip" && validation.IsOfficeDocumentContentType(provided) {
		return provided, nil
	}
	return "", fmt.Errorf("unsupported content type")
}
