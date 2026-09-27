package request

import (
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/validation"
	"net/http"
)

func GetULIDPathValue(r *http.Request, key string) (string, error) {
	value := r.PathValue(key)
	if value == "" {
		return "", fmt.Errorf("path value %q is empty: %w", key, coreerrors.ErrInvalidArgument)
	}
	if err := validation.ValidateULID(value); err != nil {
		return "", fmt.Errorf("path value %q is not a valid ULID: %w", key, coreerrors.ErrInvalidArgument)
	}
	return value, nil
}

func GetOptionalULIDQueryValue(r *http.Request, key string) (string, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return "", nil
	}
	if err := validation.ValidateULID(value); err != nil {
		return "", fmt.Errorf("query value %q is not a valid ULID: %w", key, coreerrors.ErrInvalidArgument)
	}
	return value, nil
}

func GetULIDQueryValue(r *http.Request, key string) (string, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return "", fmt.Errorf("query value %q is empty: %w", key, coreerrors.ErrInvalidArgument)
	}
	if err := validation.ValidateULID(value); err != nil {
		return "", fmt.Errorf("query value %q is not a valid ULID: %w", key, coreerrors.ErrInvalidArgument)
	}
	return value, nil
}
