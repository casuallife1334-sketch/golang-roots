package request

import (
	"encoding/json"
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
	"io"
	"net/http"
)

func DecodeJSON(r *http.Request, destination any) error {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read JSON request: %w: %w", err, coreerrors.ErrInvalidArgument)
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return fmt.Errorf("decode JSON request: %w: %w", err, coreerrors.ErrInvalidArgument)
	}
	for _, field := range []string{"birth_date", "death_date"} {
		if value, ok := payload[field].(string); ok && len(value) == len("2006-01-02") {
			payload[field] = value + "T00:00:00Z"
		}
	}
	normalized, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal JSON request: %w: %w", err, coreerrors.ErrInvalidArgument)
	}
	if err := json.Unmarshal(normalized, destination); err != nil {
		return fmt.Errorf("decode JSON request: %w: %w", err, coreerrors.ErrInvalidArgument)
	}
	return nil
}
