package storage

import (
	"context"
	"io"
)

type FileStorage interface {
	Put(context.Context, string, io.Reader, string) (string, error)
	Get(context.Context, string) (io.ReadCloser, error)
	Delete(context.Context, string) error
}
