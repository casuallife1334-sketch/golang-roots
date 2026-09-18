package service

import "context"

func (s *TreesService) DeleteTree(ctx context.Context, userID, treeID string) error {
	if err := s.CanWriteTree(ctx, userID, treeID); err != nil {
		return err
	}
	photoURLs, err := s.treesRepository.GetPhotoURLs(ctx, treeID)
	if err != nil {
		return err
	}
	for _, photoURL := range photoURLs {
		if err := s.fileStorage.Delete(ctx, photoURL); err != nil {
			return err
		}
	}
	return s.treesRepository.DeleteTree(ctx, userID, treeID)
}
