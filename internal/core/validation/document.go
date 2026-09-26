package validation

import (
	"errors"
	"github.com/oklog/ulid/v2"
	"strings"
	"unicode/utf8"
)

const MaxDocumentSize int64 = 20 << 20

var (
	ErrInvalidDocumentOwner = errors.New("invalid document owner")
	ErrInvalidDocument      = errors.New("invalid document")
	ErrDocumentTooLarge     = errors.New("document is too large")
)

var documentContentTypes = map[string]struct{}{
	"application/pdf": {},
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {},
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       {},
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
}

func ValidateDocumentOwnerType(value string) error {
	if value != "person" && value != "relationship" {
		return ErrInvalidDocumentOwner
	}
	return nil
}

func ValidateDocumentOwnerID(value string) error {
	if strings.TrimSpace(value) == "" || utf8.RuneCountInString(value) != 26 {
		return ErrInvalidDocumentOwner
	}
	if _, err := ulid.Parse(value); err != nil {
		return ErrInvalidDocumentOwner
	}
	return nil
}

func NormalizeDocumentFileName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" || utf8.RuneCountInString(value) > 255 {
		return "", ErrInvalidDocument
	}
	return value, nil
}

func ValidateDocument(contentType string, size int64) error {
	if size <= 0 {
		return ErrInvalidDocument
	}
	if size > MaxDocumentSize {
		return ErrDocumentTooLarge
	}
	if _, ok := documentContentTypes[contentType]; !ok {
		return ErrInvalidDocument
	}
	return nil
}

func IsDocumentContentType(value string) bool {
	_, ok := documentContentTypes[value]
	return ok
}

func IsOfficeDocumentContentType(value string) bool {
	return value == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" ||
		value == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
}
