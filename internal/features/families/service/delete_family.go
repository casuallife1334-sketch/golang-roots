package service

import "context"

func (s *FamiliesService) DeleteFamily(ctx context.Context, userID, treeID, id string) error {
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return err
	}
	return s.familyRepository.DeleteFamily(ctx, treeID, id)
}
