package service

import (
	"context"
	"fmt"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/validation"
	"github.com/oklog/ulid/v2"
	"path/filepath"
)

func (s *DocumentsService) CreateDocument(ctx context.Context, userID, treeID string, input UploadInput) (domain.Document, error) {
	if err := validateOwner(input.Owner); err != nil {
		return domain.Document{}, ErrInvalid
	}
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return domain.Document{}, err
	}
	fileName, err := validation.NormalizeDocumentFileName(filepath.Base(input.FileName))
	if err != nil {
		return domain.Document{}, ErrInvalid
	}
	if err := validation.ValidateDocument(input.ContentType, input.SizeBytes); err != nil {
		if err == validation.ErrDocumentTooLarge {
			return domain.Document{}, ErrTooLarge
		}
		return domain.Document{}, ErrInvalid
	}
	documentID := ulid.Make().String()
	storageKey := fmt.Sprintf("documents/%s/%s", treeID, documentID)
	if _, err := s.fileStorage.Put(ctx, storageKey, input.File, input.ContentType); err != nil {
		return domain.Document{}, fmt.Errorf("put document: %w", err)
	}
	document, err := s.documentsRepository.CreateDocument(ctx, domain.CreateDocumentInput{
		TreeID:      treeID,
		Owner:       input.Owner,
		FileName:    fileName,
		StorageKey:  storageKey,
		ContentType: input.ContentType,
		SizeBytes:   input.SizeBytes,
		CreatedBy:   userID,
	})
	if err != nil {
		_ = s.fileStorage.Delete(ctx, storageKey)
		return domain.Document{}, err
	}
	return document, nil
}

func validateOwner(owner domain.DocumentOwner) error {
	if err := validation.ValidateDocumentOwnerType(string(owner.Type)); err != nil {
		return err
	}
	if err := validation.ValidateDocumentOwnerID(owner.ID); err != nil {
		return err
	}
	return nil
}
