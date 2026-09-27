package validation

import (
	"errors"
	"github.com/oklog/ulid/v2"
)

var ErrInvalidULID = errors.New("invalid ULID")

func ValidateULID(value string) error {
	if _, err := ulid.Parse(value); err != nil {
		return ErrInvalidULID
	}
	return nil
}
