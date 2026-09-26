package service

import (
	"context"
	"fmt"
	"genealogy-tree/internal/core/domain"
	"io"
)

func (s *DocumentsService) GetDocuments(ctx context.Context, userID, treeID string, owner domain.DocumentOwner) ([]domain.Document, error) {
	if err := validateOwner(owner); err != nil {
		return nil, ErrInvalid
	}
	if err := s.treeAccess.CanReadTree(ctx, userID, treeID); err != nil {
		return nil, err
	}
	return s.documentsRepository.GetDocuments(ctx, treeID, owner)
}

func (s *DocumentsService) OpenDocument(ctx context.Context, userID, treeID, id string) (domain.Document, io.ReadCloser, error) {
	if err := s.treeAccess.CanReadTree(ctx, userID, treeID); err != nil {
		return domain.Document{}, nil, err
	}
	document, err := s.documentsRepository.GetDocument(ctx, treeID, id)
	if err != nil {
		return domain.Document{}, nil, err
	}
	file, err := s.fileStorage.Get(ctx, document.StorageKey)
	if err != nil {
		return domain.Document{}, nil, fmt.Errorf("get document: %w", err)
	}
	return document, file, nil
}
