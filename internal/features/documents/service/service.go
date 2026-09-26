package service

import (
	"context"
	"fmt"
	"genealogy-tree/internal/core/domain"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/storage"
	"io"
)

var (
	ErrInvalid  = fmt.Errorf("%w: invalid document input", coreerrors.ErrInvalidArgument)
	ErrTooLarge = fmt.Errorf("%w: document is too large", coreerrors.ErrPayloadTooLarge)
)

type DocumentsRepository interface {
	CreateDocument(context.Context, domain.CreateDocumentInput) (domain.Document, error)
	GetDocuments(context.Context, string, domain.DocumentOwner) ([]domain.Document, error)
	GetDocument(context.Context, string, string) (domain.Document, error)
	DeleteDocument(context.Context, string, string) error
}

type TreeAccess interface {
	CanReadTree(context.Context, string, string) error
	CanWriteTree(context.Context, string, string) error
}

type DocumentsService struct {
	documentsRepository DocumentsRepository
	fileStorage         storage.FileStorage
	treeAccess          TreeAccess
}

func NewDocumentsService(documentsRepository DocumentsRepository, fileStorage storage.FileStorage, treeAccess TreeAccess) *DocumentsService {
	return &DocumentsService{
		documentsRepository: documentsRepository,
		fileStorage:         fileStorage,
		treeAccess:          treeAccess,
	}
}

type UploadInput struct {
	Owner       domain.DocumentOwner
	FileName    string
	ContentType string
	SizeBytes   int64
	File        io.Reader
}
