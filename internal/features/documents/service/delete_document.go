package service

import "context"

func (s *DocumentsService) DeleteDocument(ctx context.Context, userID, treeID, id string) error {
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return err
	}
	document, err := s.documentsRepository.GetDocument(ctx, treeID, id)
	if err != nil {
		return err
	}
	if err := s.fileStorage.Delete(ctx, document.StorageKey); err != nil {
		return err
	}
	return s.documentsRepository.DeleteDocument(ctx, treeID, id)
}
