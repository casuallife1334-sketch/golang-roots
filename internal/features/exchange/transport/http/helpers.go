package http

import (
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/features/exchange/service"
	"mime/multipart"
	"net/http"
	"strings"
)

const maxGEDCOMSize = 20 << 20

func requireGEDCOMFormat(r *http.Request) error {
	if strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("format")), "gedcom") {
		return nil
	}
	return fmt.Errorf("format must be gedcom: %w", coreerrors.ErrInvalidArgument)
}

func readGEDCOMFile(w http.ResponseWriter, r *http.Request) (multipart.File, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxGEDCOMSize+(1<<20))
	if err := r.ParseMultipartForm(maxGEDCOMSize + (1 << 20)); err != nil {
		return nil, service.ErrTooLarge
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("multipart field 'file' is required: %w", coreerrors.ErrInvalidArgument)
	}
	if header.Size > maxGEDCOMSize {
		_ = file.Close()
		return nil, service.ErrTooLarge
	}
	return file, nil
}
