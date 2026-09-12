package request

import (
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
	"github.com/oklog/ulid/v2"
	"net/http"
)

func GetULIDPathValue(r *http.Request, key string) (string, error) {
	value := r.PathValue(key)
	if value == "" {
		return "", fmt.Errorf("path value %q is empty: %w", key, coreerrors.ErrInvalidArgument)
	}
	if _, err := ulid.Parse(value); err != nil {
		return "", fmt.Errorf("path value %q is not a valid ULID: %w", key, coreerrors.ErrInvalidArgument)
	}
	return value, nil
}
