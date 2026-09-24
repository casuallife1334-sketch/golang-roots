// Package nullable provides JSON values that distinguish an omitted field
// from an explicit null. It is intended for partial update contracts.
package nullable

import (
	"bytes"
	"encoding/json"
)

// Value distinguishes an omitted JSON field from an explicit null or value.
type Value[T any] struct {
	Set   bool
	Value *T
}

func (v *Value[T]) UnmarshalJSON(data []byte) error {
	v.Set = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		v.Value = nil
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	v.Value = &value
	return nil
}
