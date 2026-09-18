package service

import "context"

func (s *RelationshipsService) DeleteRelationship(ctx context.Context, userID, treeID, id string) error {
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return err
	}
	return s.relationshipsRepository.DeleteRelationship(ctx, treeID, id)
}
