package errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrConflict        = errors.New("conflict")
	ErrAmbiguousFamily = fmt.Errorf("%w: multiple families match; specify family_id", ErrConflict)
	ErrNotFound        = errors.New("not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrPayloadTooLarge = errors.New("payload too large")
)
